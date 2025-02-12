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

	"go.temporal.io/sdk/activity"
)

func CreateProjectActivity(ctx context.Context, config CloudConfig) (*ProjectResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CreateProjectActivity", config.Provider)

	switch config.Provider {
	case CloudProviderAWS:
		return createAWSProject(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

func SetupNetworkActivity(ctx context.Context, project *ProjectResult) (*NetworkResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SetupNetworkActivity")

	switch project.Provider {
	case CloudProviderAWS:
		return createAWSNetwork(ctx, project, project.Subnet)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", project.Provider)
	}
}

func ProvisionNodeActivity(ctx context.Context, project *ProjectResult) (*NodeResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProvisionNodeActivity")

	switch project.Provider {
	case CloudProviderAWS:
		return createAWSNode(ctx, project)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", project.Provider)
	}
}
