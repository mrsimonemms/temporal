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

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CloudProvisionWorkflow(ctx workflow.Context, cfg CloudConfig) (*ProjectResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting cloud provisioning workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	})

	logger.Debug("Create project in cloud provider")
	var project *ProjectResult
	if err := workflow.ExecuteActivity(ctx, CreateProjectActivity, cfg).Get(ctx, &project); err != nil {
		logger.Error("Error executing cloud provisioning activity", err)
		return nil, fmt.Errorf("error executing cloud provision activity: %w", err)
	}

	logger.Debug("Create network in cloud provider")
	var network *NetworkResult
	if err := workflow.ExecuteActivity(ctx, SetupNetworkActivity, project).Get(ctx, &network); err != nil {
		logger.Error("Error setting up network activity", err)
		return nil, fmt.Errorf("error setting up network activity: %w", err)
	}
	project.Network = network

	// Run as a child process to fan-out to support multiple node creation
	logger.Debug(("Create nodes in cloud provider"))
	project.Nodes = make([]*NodeResult, 0)

	type nodeResult struct {
		ctx workflow.Context
		w   workflow.ChildWorkflowFuture
	}

	provisionNodeFutures := []nodeResult{}

	// Invoke the child workflows in parallel
	for i := range cfg.VMCount {
		// Set ID so can track the jobs in dashboard easier
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowTaskTimeout: time.Hour,
			WorkflowID:          fmt.Sprintf("%s_node_%d", workflow.GetInfo(ctx).WorkflowExecution.ID, i),
		})

		// Execute the child workflow and store results as a Future
		provisionNodeFutures = append(provisionNodeFutures, nodeResult{
			ctx: childCtx,
			w:   workflow.ExecuteChildWorkflow(childCtx, ProvisionNodeWorkflow, project),
		})
	}

	// Now the child workflows are running, wait for the results
	for _, k := range provisionNodeFutures {
		var node *NodeResult

		if err := k.w.Get(k.ctx, &node); err != nil {
			logger.Error("Error provisioning nodes", err)
			return nil, fmt.Errorf("error provisioning nodes: %w", err)
		}

		project.Nodes = append(project.Nodes, node)
	}

	return project, nil
}

// Run as a child worker
func ProvisionNodeWorkflow(ctx workflow.Context, project *ProjectResult) (*NodeResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting node provisioning workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	})

	var node NodeResult
	err := workflow.ExecuteActivity(ctx, ProvisionNodeActivity, project).Get(ctx, &node)
	if err != nil {
		logger.Error("Error executing node provisioning activity", err)
		return nil, fmt.Errorf("error executing node provision activity: %w", err)
	}

	return &node, nil
}
