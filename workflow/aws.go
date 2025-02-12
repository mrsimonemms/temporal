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
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"go.temporal.io/sdk/activity"
)

func createAWSProject(ctx context.Context, cfg CloudConfig) (*ProjectResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate project creation job")
	time.Sleep(time.Second)

	if err := simulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	// These values may come from the project or from input variables
	return &ProjectResult{
		CloudConfig: cfg,
		ID:          uuid.NewString(),
	}, nil
}

// project isn't actually used in here - would be part of the API call to cloud in reality
func createAWSNode(ctx context.Context, _ *ProjectResult) (*NodeResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate node setup job")
	time.Sleep(time.Second * 5)

	if err := simulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	// Generate a machine name - real service could be more descriptive (pets), entirely arbitrary (cattle) or from default provider's name
	seed := time.Now().UTC().UnixNano()
	generator := namegenerator.NewNameGenerator(seed)

	return &NodeResult{
		ID:      uuid.NewString(),
		Name:    generator.Generate(),
		Address: string(generateIPAddress()),
		Port:    22,
	}, nil
}

func createAWSNetwork(ctx context.Context, cfg *ProjectResult, subnetCidr string) (*NetworkResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate network setup job")
	time.Sleep(time.Second * 5)

	if err := simulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	_, subnet, err := net.ParseCIDR(subnetCidr)
	if err != nil {
		return nil, fmt.Errorf("error parsing cidr: %w", err)
	}

	// These values may come from the project or from input variables
	return &NetworkResult{
		ID:     uuid.NewString(),
		Region: cfg.Region,
		Subnet: subnet,
	}, nil
}
