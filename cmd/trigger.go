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

package cmd

import (
	"context"
	"log/slog"

	"github.com/mrsimonemms/temporal/pkg/providers"
	"github.com/mrsimonemms/temporal/pkg/workflow"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	tLog "go.temporal.io/sdk/log"
)

var triggerOpts providers.CloudConfig

var triggerProvider string

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Run the Temporal workflow",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.Dial(client.Options{
			HostPort: rootOpts.Host,
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

		triggerOpts.Provider = providers.CloudProvider(triggerProvider)

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.CloudProvisionWorkflow, triggerOpts)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to execute workflow")
		}

		log.Info().Str("WorkflowID", we.GetID()).Str("RunID", we.GetRunID()).Msg("Started workflow")

		// Synchronously wait for the workflow completion.
		var result providers.ProjectResult
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to get workflow result")
		}
		log.Info().Interface("result", result).Msg("Workflow finished")
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)

	bindEnv("count", 3)
	triggerCmd.Flags().IntVar(&triggerOpts.VMCount, "count", viper.GetInt("count"), "Number of VMs to build")

	bindEnv("region", "eu-west-2")
	triggerCmd.Flags().StringVar(&triggerOpts.Region, "region", viper.GetString("region"), "Region in which to build the resources")

	bindEnv("subnet", "10.0.0.0/24")
	triggerCmd.Flags().StringVar(&triggerOpts.Subnet, "subnet", viper.GetString("subnet"), "Subnet to use for the network")

	bindEnv("provider", string(providers.CloudProviderAWS))
	triggerCmd.Flags().StringVar(&triggerProvider, "provider", viper.GetString("provider"), "Cloud provider to use")
}
