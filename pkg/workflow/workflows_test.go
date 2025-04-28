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

package workflow_test

import (
	"net"
	"testing"

	"github.com/mrsimonemms/temporal/pkg/providers"
	"github.com/mrsimonemms/temporal/pkg/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_CloudProvisionWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	expectedNodes := []*providers.NodeResult{
		{
			ID:      "node0",
			Name:    "node-name-0",
			Address: net.IPv4(10, 20, 30, 40),
			Port:    22,
		},
		{
			ID:      "node1",
			Name:    "node-name-1",
			Address: net.IPv4(10, 20, 30, 41),
			Port:    22,
		},
	}
	cfg := providers.CloudConfig{
		Provider: providers.CloudProviderAWS,
		VMCount:  len(expectedNodes),
	}
	expectedProject := &providers.ProjectResult{
		CloudConfig: cfg,
		ID:          "some-id",
	}
	expectedNetwork := &providers.NetworkResult{
		ID:     "some-network-id",
		Region: "some-region-id",
		// net.ParseCIDR function returns the 16-byte representation
		// @link https://github.com/golang/go/issues/35727
		Subnet: &net.IPNet{
			IP:   net.IPv4(192, 0, 2, 0),
			Mask: net.IPv4Mask(192, 0, 2, 0),
		},
	}

	// Mock the activity responses
	env.OnActivity(workflow.CreateProjectActivity, mock.Anything, cfg).Return(expectedProject, nil)
	env.OnActivity(workflow.SetupNetworkActivity, mock.Anything, cfg, expectedProject).Return(expectedNetwork, nil)

	// Mock the child workflow
	env.RegisterWorkflow(workflow.ProvisionNodeWorkflow)
	for _, node := range expectedNodes {
		env.OnWorkflow("ProvisionNodeWorkflow", mock.Anything, mock.Anything, mock.Anything).Return(node, nil).Once()
	}

	env.ExecuteWorkflow(workflow.CloudProvisionWorkflow, cfg)
	assert.True(t, env.IsWorkflowCompleted())

	var result *providers.ProjectResult
	assert.NoError(t, env.GetWorkflowResult(&result))

	assert.Equal(t, expectedProject.CloudConfig, result.CloudConfig)
	assert.Equal(t, expectedProject.ID, result.ID)
	assert.Equal(t, expectedNetwork, result.Network)
	assert.ElementsMatch(t, expectedNodes, result.Nodes)

	env.AssertExpectations(t)
}

func Test_ProvisionNodeWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	expectedNode := &providers.NodeResult{
		ID:      "some-id",
		Name:    "some name",
		Address: net.IPv4(10, 20, 30, 40),
		Port:    22,
	}
	expectedNodeReady := &providers.NodeReadyResult{
		Ready: true,
	}

	cfg := providers.CloudConfig{
		Provider: providers.CloudProviderAWS,
	}
	project := &providers.ProjectResult{
		CloudConfig: cfg,
		Nodes:       []*providers.NodeResult{},
	}

	// Mock the activity responses
	env.OnActivity(workflow.ProvisionNodeActivity, mock.Anything, cfg, project).Return(expectedNode, nil)
	env.OnActivity(workflow.AwaitForNodeRunningActivity, mock.Anything, cfg, expectedNode).Return(expectedNodeReady, nil)

	env.ExecuteWorkflow(workflow.ProvisionNodeWorkflow, cfg, project)
	assert.True(t, env.IsWorkflowCompleted())

	var result *providers.NodeResult
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, expectedNode, result)

	env.AssertExpectations(t)
}
