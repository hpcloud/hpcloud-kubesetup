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

package blockstorage_test

import (
	"testing"

	blockstorage "git.openstack.org/stackforge/golang-client.git/blockstorage/v1"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/misc"
)

func TestVolumeTypeScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	blockstorageService := blockstorage.NewService(authenticator)

	volumeTypes, err := blockstorageService.VolumeTypes()
	if err != nil {
		t.Fatal("Cannot access volumeTypes:", err)
	}

	if len(volumeTypes) > 0 {
		_, err := blockstorageService.VolumeType(volumeTypes[0].ID)
		if err != nil {
			t.Fatal("Cannot requery volume type:", err)
		}
	}

	createVolumeParameters := blockstorage.CreateVolumeTypeParameters{Name: "nonstandard"}

	createdVolumeType, err := blockstorageService.CreateVolumeType(createVolumeParameters)

	skipRestOfTest := false
	if err != nil {
		status, ok := err.(misc.HTTPStatus)

		if ok && status.StatusCode == 404 || status.StatusCode == 403 {
			t.Log("Cannot test volume type create, so not attempting to delete.")
			skipRestOfTest = true
		} else {
			t.Fatal("Cannot create volume type:", err)
		}
	}

	if !skipRestOfTest {
		queriedVolumeType, err := blockstorageService.VolumeType(createdVolumeType.ID)
		if err != nil {
			t.Fatal("Cannot requery volume type:", err)
		}

		err = blockstorageService.DeleteVolumeType(queriedVolumeType.ID)
		if err != nil {
			t.Fatal("Cannot delete volume type:", err)
		}
	}
}
