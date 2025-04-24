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

package providers

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"go.temporal.io/sdk/activity"
)

type aws struct {
	cfg *CloudConfig
}

// CheckNodeReady implements Provider.
func (a aws) CheckNodeReady(ctx context.Context, node *NodeResult) error {
	// Generate a local timeout - this is not a Temporal sleep, but exists to
	// simulate the time taken by the VM's SSH server to become ready.
	minValue := 1
	maxValue := 30
	//nolint:gosec // ignore weak number generator error
	timeoutLength := rand.IntN(maxValue-minValue+1) + minValue
	timeout := time.Duration(timeoutLength) * time.Second

	logger := activity.GetLogger(ctx)
	logger.Info("Timing out", "timeout", timeout)

	time.Sleep(timeout)
	return nil
}

func (a aws) CreateNetwork(ctx context.Context, project *ProjectResult) (*NetworkResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate network setup job")
	time.Sleep(time.Second * 5)

	if err := SimulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	_, subnet, err := net.ParseCIDR(project.Subnet)
	if err != nil {
		return nil, fmt.Errorf("error parsing cidr: %w", err)
	}

	// These values may come from the project or from input variables
	return &NetworkResult{
		ID:     uuid.NewString(),
		Region: a.cfg.Region,
		Subnet: subnet,
	}, nil
}

func (a aws) CreateNode(ctx context.Context, project *ProjectResult) (*NodeResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate node setup job")
	time.Sleep(time.Second * 5)

	if err := SimulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	// Generate a machine name - real service could be more descriptive (pets), entirely arbitrary (cattle) or from default provider's name
	seed := time.Now().UTC().UnixNano()
	generator := namegenerator.NewNameGenerator(seed)

	return &NodeResult{
		ID:      uuid.NewString(),
		Name:    generator.Generate(),
		Address: GenerateIPAddress(),
		Port:    22,
	}, nil
}

func (a aws) CreateProject(ctx context.Context) (*ProjectResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Debug("Sleeping to simulate project creation job")
	time.Sleep(time.Second)

	if err := SimulateFailure(); err != nil {
		return nil, fmt.Errorf("simulated cloud failure: %w", err)
	}

	// These values may come from the project or from input variables
	return &ProjectResult{
		CloudConfig: *a.cfg,
		ID:          uuid.NewString(),
	}, nil
}

func NewAWS(cfg *CloudConfig) (Provider, error) {
	return aws{
		cfg: cfg,
	}, nil
}
