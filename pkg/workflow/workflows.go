/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package workflow

import (
	"fmt"
	"time"

	"github.com/mrsimonemms/temporal/pkg/providers"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CloudProvisionWorkflow(ctx workflow.Context, cfg providers.CloudConfig) (*providers.ProjectResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting cloud provisioning workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
		},
	})

	logger.Debug("Create project in cloud provider")
	var project *providers.ProjectResult
	if err := workflow.ExecuteActivity(ctx, CreateProjectActivity, cfg).Get(ctx, &project); err != nil {
		logger.Error("Error executing cloud provisioning activity", "error", err)
		return nil, fmt.Errorf("error executing cloud provision activity: %w", err)
	}

	logger.Debug("Create network in cloud provider")
	var network *providers.NetworkResult
	if err := workflow.ExecuteActivity(ctx, SetupNetworkActivity, cfg, project).Get(ctx, &network); err != nil {
		logger.Error("Error setting up network activity", "error", err)
		return nil, fmt.Errorf("error setting up network activity: %w", err)
	}
	project.Network = network

	// Run as a child process to fan-out to support multiple node creation
	logger.Debug("Create nodes in cloud provider")
	project.Nodes = make([]*providers.NodeResult, 0)

	provisionNodeFutures := map[workflow.Context]workflow.ChildWorkflowFuture{}

	// Invoke the child workflows in parallel
	for i := range cfg.VMCount {
		// Set ID so can track the jobs in dashboard easier
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowTaskTimeout: time.Hour,
			WorkflowID:          fmt.Sprintf("%s_node_%d", workflow.GetInfo(ctx).WorkflowExecution.ID, i),
		})

		// Execute the child workflow and store results as a Future
		provisionNodeFutures[childCtx] = workflow.ExecuteChildWorkflow(childCtx, ProvisionNodeWorkflow, cfg, project)
	}

	// Now the child workflows are running, wait for the results
	for ctx, workflow := range provisionNodeFutures {
		var node *providers.NodeResult

		if err := workflow.Get(ctx, &node); err != nil {
			logger.Error("Error provisioning nodes", "error", err)
			return nil, fmt.Errorf("error provisioning nodes: %w", err)
		}

		project.Nodes = append(project.Nodes, node)
	}

	return project, nil
}

// Run as a child worker
func ProvisionNodeWorkflow(
	ctx workflow.Context,
	cfg providers.CloudConfig,
	project *providers.ProjectResult,
) (*providers.NodeResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting node provisioning workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
		},
	})

	var node *providers.NodeResult
	if err := workflow.ExecuteActivity(ctx, ProvisionNodeActivity, cfg, project).Get(ctx, &node); err != nil {
		logger.Error("Error executing node provisioning activity", "error", err)
		return nil, fmt.Errorf("error executing node provision activity: %w", err)
	}

	var isReady *providers.NodeReadyResult
	if err := workflow.ExecuteActivity(ctx, AwaitForNodeRunningActivity, cfg, node).Get(ctx, &isReady); err != nil {
		logger.Error("Error whilst waiting for node to become ready", "error", err)
		return nil, fmt.Errorf("error waiting for node to become ready: %w", err)
	}

	return node, nil
}
