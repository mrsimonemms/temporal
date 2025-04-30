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

package temporal

import (
	"context"
	"crypto/tls"
	"log/slog"

	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"go.temporal.io/sdk/client"
	tLog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func apiKeyInterceptor(namespace string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(
			metadata.AppendToOutgoingContext(ctx, "temporal-namespace", namespace),
			method,
			req,
			reply,
			cc,
			opts...,
		)
	}
}

// Common function to create a Temporal client
func NewClient(host, namespace, apiKey string) (client.Client, error) {
	var credentials client.Credentials
	var connectionOptions client.ConnectionOptions

	// These are the changes required to use Temporal Cloud with an API key
	if apiKey != "" {
		credentials = client.NewAPIKeyStaticCredentials(apiKey)

		connectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(apiKeyInterceptor(namespace)),
			},
		}
	}
	return client.Dial(client.Options{
		HostPort:          host,
		Namespace:         namespace,
		ConnectionOptions: connectionOptions,
		Credentials:       credentials,
		DataConverter:     NewDataConverter(),
		FailureConverter: temporal.NewDefaultFailureConverter(temporal.DefaultFailureConverterOptions{
			EncodeCommonAttributes: true,
		}),
		Logger: tLog.NewStructuredLogger(slog.New(slogzerolog.Option{
			Logger: &log.Logger,
		}.NewZerologHandler())),
	})
}
