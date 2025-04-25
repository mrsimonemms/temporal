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
	"context"
	"fmt"

	"github.com/mrsimonemms/temporal/pkg/providers"
	"go.temporal.io/sdk/activity"
)

func CreateProjectActivity(ctx context.Context, config providers.CloudConfig) (*providers.ProjectResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CreateProjectActivity", "provider", config.Provider)

	cloudProvider, err := config.GetProvider()
	if err != nil {
		return nil, fmt.Errorf("error initializing provider: %w", err)
	}

	return cloudProvider.CreateProject(ctx)
}

func SetupNetworkActivity(
	ctx context.Context,
	config providers.CloudConfig,
	project *providers.ProjectResult,
) (*providers.NetworkResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SetupNetworkActivity", "provider", config.Provider)

	cloudProvider, err := config.GetProvider()
	if err != nil {
		return nil, fmt.Errorf("error initializing provider: %w", err)
	}

	return cloudProvider.CreateNetwork(ctx, project)
}

// Simulate making an SSH connection and checking for cloud-config to become ready
func AwaitForNodeRunningActivity(
	ctx context.Context,
	config providers.CloudConfig,
	node *providers.NodeResult,
) (*providers.NodeReadyResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProvisionNodeActivity", "provider", config.Provider)

	cloudProvider, err := config.GetProvider()
	if err != nil {
		return nil, fmt.Errorf("error initializing provider: %w", err)
	}

	if err := cloudProvider.CheckNodeReady(ctx, node); err != nil {
		return nil, err
	}

	// If there's no error then it's ready
	return &providers.NodeReadyResult{Ready: true}, nil
}

func ProvisionNodeActivity(ctx context.Context,
	config providers.CloudConfig,
	project *providers.ProjectResult,
) (*providers.NodeResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProvisionNodeActivity", "provider", config.Provider)

	cloudProvider, err := config.GetProvider()
	if err != nil {
		return nil, fmt.Errorf("error initializing provider: %w", err)
	}

	return cloudProvider.CreateNode(ctx, project)
}
