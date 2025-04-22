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
	"log/slog"
	"os"

	"github.com/mrsimonemms/temporal/workflow"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"go.temporal.io/sdk/client"
	tLog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
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

	w := worker.New(c, "cloud-provisioning", worker.Options{})

	// Register the workflows
	w.RegisterWorkflow(workflow.ProvisionNodeWorkflow)
	w.RegisterWorkflow(workflow.CloudProvisionWorkflow)

	// Register the activities
	w.RegisterActivity(workflow.CreateProjectActivity)
	w.RegisterActivity(workflow.SetupNetworkActivity)
	w.RegisterActivity(workflow.ProvisionNodeActivity)
	w.RegisterActivity(workflow.AwaitForNodeRunningActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start worker")
	}
}
