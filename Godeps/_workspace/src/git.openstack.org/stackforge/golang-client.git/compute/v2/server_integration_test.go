// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package compute_test

import (
	"os"
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
)

func TestServerScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	imageID := os.Getenv("TEST_IMAGEID")
	if imageID == "" {
		t.Skip("Cannot run test because Env Var 'TEST_IMAGEID' is not set")
	}

	keyPairName := os.Getenv("TEST_KEYPAIR_NAME")
	if keyPairName == "" {
		t.Skip("Cannot run test because Env Var 'TEST_KEYPAIR_NAME' is not set")
	}

	flavorRef := os.Getenv("TEST_FLAVOR")
	if flavorRef == "" {
		t.Skip("Cannot run test because Env Var 'TEST_FLAVOR' is not set")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	computeService := compute.NewService(authenticator)

	// Create a server with the min required parameters
	minParametersServer := compute.ServerCreationParameters{
		Name:        "testName",
		ImageRef:    imageID,
		KeyPairName: keyPairName,
		FlavorRef:   flavorRef,
	}

	createQueryDeleteServer(t, computeService, minParametersServer)
}

func createQueryDeleteServer(t *testing.T, computeService compute.Service, scp compute.ServerCreationParameters) {
	createdServer, err := computeService.CreateServer(scp)
	if err != nil {
		t.Fatal("Cannot create server:", err)
	}

	queriedServer, err := computeService.ServerDetail(createdServer.ID)
	if err != nil {
		t.Fatal("Cannot requery server:", err)
	}

	servers, err := computeService.Servers()
	if err != nil {
		t.Fatal("Cannot access Servers:", err)
	}

	foundServer := false
	for _, serverValue := range servers {
		if queriedServer.ID == serverValue.ID {
			foundServer = true
			break
		}
	}

	if !foundServer {
		t.Fatal("Cannot find server that was created.")
	}

	err = computeService.DeleteServer(queriedServer.ID)
	if err != nil {
		t.Fatal("Cannot delete the server:", err)
	}
}
