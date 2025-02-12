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
	"os"

	"github.com/mrsimonemms/temporal/workflows/helloworld"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func main() {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = client.DefaultHostPort
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create Temporal client")
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, "Simon")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to execute workflow")
	}

	log.Info().Str("WorkflowID", we.GetID()).Str("RunID", we.GetRunID()).Msg("Started workflow")

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get workflow result")
	}
	log.Info().Str("result", result).Msg("Workflow finished")
}
