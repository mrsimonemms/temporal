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
	"os"
	"strings"

	"github.com/mrsimonemms/temporal/pkg/temporal"
	"github.com/mrsimonemms/temporal/pkg/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var rootOpts struct {
	APIKey    string
	Host      string
	Namespace string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "temporal",
	Short: "Temporal demo application",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := temporal.NewClient(rootOpts.Host, rootOpts.Namespace, rootOpts.APIKey)
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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func bindEnv(key string, defaultValue ...any) {
	envvarName := strings.Replace(key, "-", "_", -1)
	envvarName = strings.ToUpper(envvarName)

	err := viper.BindEnv(key, envvarName)
	cobra.CheckErr(err)

	for _, val := range defaultValue {
		viper.SetDefault(key, val)
	}
}

func init() {
	bindEnv("temporal-address", client.DefaultHostPort)
	rootCmd.PersistentFlags().StringVarP(&rootOpts.Host, "temporal-address", "a", viper.GetString("temporal-address"), "Address for Temporal server")

	bindEnv("temporal-key", "")
	rootCmd.PersistentFlags().StringVarP(&rootOpts.APIKey, "temporal-key", "k", viper.GetString("temporal-key"), "API key for Temporal server")
	rootCmd.PersistentFlags().Lookup("temporal-key").DefValue = "xxx"

	bindEnv("temporal-namespace", client.DefaultNamespace)
	rootCmd.PersistentFlags().StringVarP(
		&rootOpts.Namespace,
		"temporal-namespace",
		"n",
		viper.GetString("temporal-namespace"),
		"Namespace for Temporal server",
	)
}
