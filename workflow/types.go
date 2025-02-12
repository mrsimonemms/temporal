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

package workflow

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
)

type ProjectResult struct {
	CloudConfig

	ID string

	Network *NetworkResult
	Nodes   []*NodeResult
}

type NetworkResult struct {
	ID     string
	Region string
	Subnet *net.IPNet
}

type NodeResult struct {
	ID      string
	Name    string
	Address string
	Port    int32
}

type Resource struct {
	ID string
}

type CloudProvider string

const (
	CloudProviderAWS   CloudProvider = "AWS"
	CloudProviderAzure CloudProvider = "Azure"
)

type CloudConfig struct {
	Provider CloudProvider
	Region   string
	Subnet   string
	VMCount  int
}

// Pseudo-randomise failure - this is obviously not going to be in a real-world
// version, but it exists to demonstrate that cloud APIs are a black box and we
// have no control over the failures
func simulateFailure() error {
	if rand.IntN(9) == 1 {
		return fmt.Errorf("simulate failure")
	}
	return nil
}

// Generate an IP address - this simulates the cloud provider's process of assigning an IP
func generateIPAddress() []byte {
	buf := make([]byte, 4)

	ip := rand.Uint32()

	binary.LittleEndian.PutUint32(buf, ip)
	return net.IP(buf)
}
