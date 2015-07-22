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
	"os"
	"testing"

	blockstorage "git.openstack.org/stackforge/golang-client.git/blockstorage/v1"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
)

func TestVolumeScenarios(t *testing.T) {
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

	volumes, err := blockstorageService.Volumes()
	if err != nil {
		t.Fatal("Cannot access volumes:", err)
	}

	if len(volumes) > 0 {
		_, err := blockstorageService.Volume(volumes[0].ID)
		if err != nil {
			t.Fatal("Cannot requery volume:", err)
		}
	}

	createVolumeParameters := blockstorage.CreateVolumeParameters{
		DisplayName: "TestVolumeDisplayName",
		Desciption:  "test description",
		VolumeType:  volumeTypes[0].ID,
		Size:        1,
	}

	createdVolume, err := blockstorageService.CreateVolume(createVolumeParameters)
	if err != nil {
		t.Fatal("Cannot create volume:", err)
	}

	queriedVolume, err := blockstorageService.Volume(createdVolume.ID)
	if err != nil {
		t.Fatal("Cannot requery volume:", err)
	}

	err = blockstorageService.DeleteVolume(queriedVolume.ID)
	if err != nil {
		t.Fatal("Cannot delete volume:", err)
	}
}

func TestVolumeFromImageScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	imageID := os.Getenv("VATestImageID")

	if imageID == "" {
		t.Skip("Skipping test as not all env variables are set: 'VATestImageID'")
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

	volumes, err := blockstorageService.Volumes()
	if err != nil {
		t.Fatal("Cannot access volumes:", err)
	}

	if len(volumes) > 0 {
		_, err := blockstorageService.Volume(volumes[0].ID)
		if err != nil {
			t.Fatal("Cannot requery volume:", err)
		}
	}

	createVolumeParameters := blockstorage.CreateVolumeParameters{
		DisplayName: "TestVolumeDisplayNameFromImage",
		Desciption:  "test description",
		VolumeType:  volumeTypes[0].ID,
		Size:        5,
		ImageRef:    imageID,
	}

	createdVolume, err := blockstorageService.CreateVolume(createVolumeParameters)
	if err != nil {
		t.Fatal("Cannot create volume:", err)
	}

	queriedVolume, err := blockstorageService.Volume(createdVolume.ID)
	if err != nil {
		t.Fatal("Cannot requery volume:", err)
	}

	err = blockstorageService.DeleteVolume(queriedVolume.ID)
	if err != nil {
		t.Fatal("Cannot delete volume:", err)
	}
}
