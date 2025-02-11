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
	"fmt"
	"log/slog"

	"github.com/mrsimonemms/temporal/pkg/workflows/helloworld"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

var workerOpts struct {
	TemporalAddress string
}

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		zlog.Debug().Msg("Connecting to temporal server")
		c, err := client.Dial(client.Options{
			HostPort: appOpts.TemporalAddress,
			Logger: log.NewStructuredLogger(slog.New(slogzerolog.Option{
				Logger: &zerolog.Logger{},
			}.NewZerologHandler())),
		})
		if err != nil {
			return fmt.Errorf("error connecting to temporal: %w", err)
		}
		defer c.Close()

		w := worker.New(c, "hello-world", worker.Options{})

		w.RegisterWorkflow(helloworld.Workflow)
		w.RegisterActivity(helloworld.Activity)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			return fmt.Errorf("unable to start worker: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)

	workerCmd.Flags().StringVar(&workerOpts.TemporalAddress, "temporal-address", bindEnv[string]("temporal-address", "localhost:7233"), "Help message for toggle")
}
