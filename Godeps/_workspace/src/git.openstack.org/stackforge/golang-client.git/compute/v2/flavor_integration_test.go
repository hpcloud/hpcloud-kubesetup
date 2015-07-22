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
	"fmt"
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestFlavorScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	computeService := compute.NewService(authenticator)

	flavors, err := computeService.Flavors()
	if err != nil {
		t.Fatalf("Cannot get flavors: %v", err)
	}

	flavorsDetail, err := computeService.FlavorsDetail()
	if err != nil {
		t.Fatalf("Cannot query FlavorsDetail: %v", err)
	}

	testUtil.Equals(t, len(flavorsDetail), len(flavors))

	if len(flavors) > 0 {
		fmt.Println("FlavorID to query:", flavors[0].ID)
		queriedItem, err := computeService.FlavorDetail(flavors[0].ID)
		fmt.Println("Results1:", queriedItem, "Error:", err)
		fmt.Println("Expected:", flavors[0])
		if err != nil {
			t.Fatalf("Cannot requery single flavor: %v", err)
		}

		testUtil.Equals(t, flavorsDetail[0], queriedItem)
	}

}
