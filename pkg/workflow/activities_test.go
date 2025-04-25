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

//nolint:dupl

package workflow_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrsimonemms/temporal/pkg/providers"
	"github.com/mrsimonemms/temporal/pkg/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

type MockedProvider struct {
	mock.Mock
}

func (m *MockedProvider) CheckNodeReady(ctx context.Context, node *providers.NodeResult) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedProvider) CreateNetwork(ctx context.Context, project *providers.ProjectResult) (*providers.NetworkResult, error) {
	args := m.Called()
	return args.Get(0).(*providers.NetworkResult), args.Error(1)
}

func (m *MockedProvider) CreateNode(ctx context.Context, project *providers.ProjectResult) (*providers.NodeResult, error) {
	args := m.Called()
	return args.Get(0).(*providers.NodeResult), args.Error(1)
}

func (m *MockedProvider) CreateProject(ctx context.Context) (*providers.ProjectResult, error) {
	args := m.Called()
	return args.Get(0).(*providers.ProjectResult), args.Error(1)
}

func Test_CreateProjectActivity(t *testing.T) {
	tests := []struct {
		Name   string
		Result *providers.ProjectResult
		Err    error
	}{
		{
			Name:   "valid provider",
			Result: &providers.ProjectResult{},
		},
		{
			Name: "invalid provider",
			Err:  fmt.Errorf("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)

			// Create the test suite
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(workflow.CreateProjectActivity)

			config := providers.CloudConfig{}

			// Create a mocked provider
			mockedProvider := new(MockedProvider)

			// Mock the GetProvider function and restore after run
			orig := providers.GetProvider
			defer func() {
				providers.GetProvider = orig
			}()
			providers.GetProvider = func(c providers.CloudConfig) (providers.Provider, error) {
				return mockedProvider, test.Err
			}

			mockedProvider.On("CreateProject").Return(test.Result, nil)

			val, err := env.ExecuteActivity(workflow.CreateProjectActivity, config)

			if test.Result != nil {
				assert.NoError(err)

				var project *providers.ProjectResult
				assert.NoError(val.Get(&project))
				assert.Equal(project, test.Result)

				mockedProvider.AssertExpectations(t)
				mockedProvider.AssertCalled(t, "CreateProject")
			}

			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
			}
		})
	}
}

func Test_SetupNetworkActivity(t *testing.T) {
	tests := []struct {
		Name   string
		Result *providers.NetworkResult
		Err    error
	}{
		{
			Name:   "valid provider",
			Result: &providers.NetworkResult{},
		},
		{
			Name: "invalid provider",
			Err:  fmt.Errorf("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)

			// Create the test suite
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(workflow.SetupNetworkActivity)

			config := providers.CloudConfig{}
			project := &providers.ProjectResult{}

			// Create a mocked provider
			mockedProvider := new(MockedProvider)

			// Mock the GetProvider function and restore after run
			orig := providers.GetProvider
			defer func() {
				providers.GetProvider = orig
			}()
			providers.GetProvider = func(c providers.CloudConfig) (providers.Provider, error) {
				return mockedProvider, test.Err
			}

			mockedProvider.On("CreateNetwork").Return(test.Result, nil)

			val, err := env.ExecuteActivity(workflow.SetupNetworkActivity, config, project)

			if test.Result != nil {
				assert.NoError(err)

				var network *providers.NetworkResult
				assert.NoError(val.Get(&network))
				assert.Equal(network, test.Result)

				mockedProvider.AssertExpectations(t)
				mockedProvider.AssertCalled(t, "CreateNetwork")
			}

			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
			}
		})
	}
}

func Test_AwaitForNodeRunningActivity(t *testing.T) {
	tests := []struct {
		Name       string
		NodeResult *providers.NodeResult
		Err        error
	}{
		{
			Name:       "valid provider",
			NodeResult: &providers.NodeResult{},
		},
		{
			Name: "invalid provider",
			Err:  fmt.Errorf("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)

			// Create the test suite
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(workflow.AwaitForNodeRunningActivity)

			config := providers.CloudConfig{}

			// Create a mocked provider
			mockedProvider := new(MockedProvider)

			// Mock the GetProvider function and restore after run
			orig := providers.GetProvider
			defer func() {
				providers.GetProvider = orig
			}()
			providers.GetProvider = func(c providers.CloudConfig) (providers.Provider, error) {
				return mockedProvider, test.Err
			}

			mockedProvider.On("CheckNodeReady").Return(nil)

			val, err := env.ExecuteActivity(workflow.AwaitForNodeRunningActivity, config, test.NodeResult)

			if test.NodeResult != nil {
				assert.NoError(err)
				fmt.Println(val)

				var nodeReady *providers.NodeReadyResult
				assert.NoError(val.Get(&nodeReady))
				assert.Equal(nodeReady, &providers.NodeReadyResult{Ready: true})

				mockedProvider.AssertExpectations(t)
				mockedProvider.AssertCalled(t, "CheckNodeReady")
			}

			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
			}
		})
	}
}

func Test_ProvisionNodeActivity(t *testing.T) {
	tests := []struct {
		Name   string
		Result *providers.NodeResult
		Err    error
	}{
		{
			Name:   "valid provider",
			Result: &providers.NodeResult{},
		},
		{
			Name: "invalid provider",
			Err:  fmt.Errorf("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)

			// Create the test suite
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestActivityEnvironment()
			env.RegisterActivity(workflow.ProvisionNodeActivity)

			config := providers.CloudConfig{}
			project := &providers.ProjectResult{}

			// Create a mocked provider
			mockedProvider := new(MockedProvider)

			// Mock the GetProvider function and restore after run
			orig := providers.GetProvider
			defer func() {
				providers.GetProvider = orig
			}()
			providers.GetProvider = func(c providers.CloudConfig) (providers.Provider, error) {
				return mockedProvider, test.Err
			}

			mockedProvider.On("CreateNode").Return(test.Result, nil)

			val, err := env.ExecuteActivity(workflow.ProvisionNodeActivity, config, project)

			if test.Result != nil {
				assert.NoError(err)

				var node *providers.NodeResult
				assert.NoError(val.Get(&node))
				assert.Equal(node, test.Result)

				mockedProvider.AssertExpectations(t)
				mockedProvider.AssertCalled(t, "CreateNode")
			}

			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
			}
		})
	}
}
