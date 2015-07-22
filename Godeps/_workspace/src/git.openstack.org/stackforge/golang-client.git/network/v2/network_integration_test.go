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

package network_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	common "git.openstack.org/stackforge/golang-client.git/identity/common"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
	network "git.openstack.org/stackforge/golang-client.git/network/v2"
)

func TestNetworkServiceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	if os.Getenv("OS_AUTH_URL") == "" {
		t.Skip("No openstack auth configured so skipping test")
	}

	authParameters, err := common.FromEnvVars()
	if err != nil {
		t.Fatal("There was an error getting authentication env variables:", err)
	}

	authenticator := identity.Authenticate(authParameters)
	authenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))

	networkService := network.NewService(&authenticator)

	networkName := "OtherNetwork"

	createdNetwork, err := networkService.CreateNetwork(
		network.CreateNetworkParameters{
			AdminStateUp: true,
			Shared:       false,
			Name:         networkName,
			TenantID:     authParameters.TenantID})
	if err != nil {
		t.Fatal("There was an error creating a network:", err)
	}

	networks, err := networkService.Networks()
	if err != nil {
		t.Fatal("There was an error getting a list of networks:", err)
	}

	networkIDs, err := networkService.NetworkIDsByName(networkName)
	if err != nil {
		t.Fatal("There was an error getting a list of networks ids by name", err)
	}

	if len(networkIDs) > 1 {
		t.Fatalf("Wrong number of networks with the name %s: %v", networkName, err)
	}

	if networkIDs[0] != createdNetwork.ID {
		t.Fatalf("Didn't find the correct networkID by its name %s: %v", networkName, err)
	}

	foundCreatedNetwork := false
	for _, networkFound := range networks {
		if reflect.DeepEqual(createdNetwork, networkFound) {
			foundCreatedNetwork = true
		}
	}

	if !foundCreatedNetwork {
		t.Fatal("Cannot find network called OtherNetwork when getting a list of networks.")
	}

	// Might be nice to have some sugar api that can do this easily for a developer...
	// Keep iterating until active or until more than 10 tries has been exceeded.
	numTries := 0
	activeNetwork := createdNetwork
	for numTries < 10 || activeNetwork.Status != "ACTIVE" {
		activeNetwork, _ = networkService.Network(createdNetwork.ID)
		numTries++
		fmt.Println("Sleeping 50ms on try:" + string(numTries) + " with status currently " + activeNetwork.Status)
		sleepDuration, _ := time.ParseDuration("50ms")
		time.Sleep(sleepDuration)
	}

	foundNetworksByName, err := networkService.QueryNetworks(network.QueryParameters{
		Name: networkName,
	})

	if err != nil {
		t.Fatal("There was an error getting a list of networks by name", err)
	}

	if len(foundNetworksByName) > 1 {
		t.Fatalf("Wrong number of networks with the name %s: %v", networkName, err)
	}

	if foundNetworksByName[0].Name != networkName && foundNetworksByName[0].ID != createdNetwork.ID {
		t.Fatalf("Didn't find the correct network by its name %s: %v", networkName, err)
	}

	err = networkService.DeleteNetwork(activeNetwork.ID)
	if err != nil {
		t.Fatal("Error in deleting OtherNetwork:", err)
	}

	networks, err = networkService.Networks()
	if err != nil {
		t.Fatal("There was an error getting a list of networks to verify the network was deleted:", err)
	}

	networkDeleted := true
	for _, networkFound := range networks {
		if reflect.DeepEqual(createdNetwork, networkFound) {
			networkDeleted = false
		}
	}

	if !networkDeleted {
		t.Fatal("Delete of 'NewNetwork' did not occur")
	}
}
