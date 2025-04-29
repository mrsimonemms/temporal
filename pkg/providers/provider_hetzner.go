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
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	hga "github.com/mrsimonemms/hetzner-golang-actions"
)

type hetzner struct {
	cfg    *CloudConfig
	client *hcloud.Client
}

func (h hetzner) CheckNodeReady(ctx context.Context, node *NodeResult) error {
	panic("unimplemented")
}

func (h hetzner) CreateNetwork(ctx context.Context, project *ProjectResult) (*NetworkResult, error) {
	panic("unimplemented")
}

func (h hetzner) CreateNode(ctx context.Context, project *ProjectResult) (*NodeResult, error) {
	server, _, err := h.client.Server.Create(ctx, hcloud.ServerCreateOpts{})
	if err != nil {
		return nil, fmt.Errorf("error creating hetzner server: %w", err)
	}

	result := &NodeResult{
		ID:      strconv.FormatInt(server.Server.ID, 10),
		Name:    server.Server.Name,
		Address: server.Server.PublicNet.IPv4.IP,
		Port:    22,
	}

	if err := hga.NewWaiter(h.client, hga.WithTimeout(time.Minute*5)).
		Wait(ctx, server.Action, server.NextActions...); err != nil {
		// Return the node so it can be deleted
		return result, fmt.Errorf("node not built within timeout: %w", err)
	}

	return result, nil
}

// Hetzner doesn't have the concept of a project as that must created in the UI to get the API token
func (h hetzner) CreateProject(ctx context.Context) (*ProjectResult, error) {
	return &ProjectResult{
		CloudConfig: *h.cfg,
	}, nil
}

func NewHetzner(cfg *CloudConfig) (Provider, error) {
	return hetzner{
		cfg:    cfg,
		client: hcloud.NewClient(),
	}, nil
}
