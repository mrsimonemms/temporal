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

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mrsimonemms/temporal/workflow"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"go.temporal.io/sdk/client"
	tLog "go.temporal.io/sdk/log"
)

func main() {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = client.DefaultHostPort
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
		Logger: tLog.NewStructuredLogger(slog.New(slogzerolog.Option{
			Logger: &log.Logger,
		}.NewZerologHandler())),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create Temporal client")
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "cloud-provisioning",
	}

	// Define the inputs in this starter - this could be set via envvars or gRPC for production
	config := workflow.CloudConfig{
		Provider: workflow.CloudProviderAWS,
		Region:   "eu-west-2",
		Subnet:   "10.0.0.0/24",
		VMCount:  3,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.CloudProvisionWorkflow, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to execute workflow")
	}

	log.Info().Str("WorkflowID", we.GetID()).Str("RunID", we.GetRunID()).Msg("Started workflow")

	// Synchronously wait for the workflow completion.
	var result workflow.ProjectResult
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get workflow result")
	}
	log.Info().Interface("result", result).Msg("Workflow finished")
}
