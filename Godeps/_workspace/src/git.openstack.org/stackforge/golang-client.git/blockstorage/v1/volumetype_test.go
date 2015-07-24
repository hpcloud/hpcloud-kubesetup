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
	"encoding/json"
	"fmt"
	"testing"

	blockstorage "git.openstack.org/stackforge/golang-client.git/blockstorage/v1"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var volumeTypeSample = blockstorage.VolumeType{
	Name:       "Display",
	ID:         "ID",
	ExtraSpecs: blockstorage.VolumeTypeExtraSpecs{Capabilities: "gpu"},
}

var volumeTypeJson, _ = json.Marshal(volumeTypeSample)

var sampleVolumeTypesJson = fmt.Sprintf(`{"volume_types" : [ %s ] }`, string(volumeTypeJson))

var sampleVolumeTypeJson = fmt.Sprintf(`{"volume_type" : %s }`, string(volumeTypeJson))

func TestGetVolumeTypes(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyStatusAndURL(t, tokn, 200, sampleVolumeTypesJson, "/types")
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	volumes, err := service.VolumeTypes()
	testUtil.IsNil(t, err)
	testUtil.Assert(t, len(volumes) == 1, "Expected 1 volume")
	testUtil.Equals(t, volumeTypeSample, volumes[0])
}

func TestGetVolumeType(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyStatusAndURL(t, tokn, 200, sampleVolumeTypeJson, "/types/id")
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	volume, err := service.VolumeType("id")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, volumeTypeSample, volume)
}

func TestDeleteVolumeType(t *testing.T) {
	name := "id"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "types/"+name)
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	err := service.DeleteVolumeType(name)
	testUtil.IsNil(t, err)
}

func TestCreateVolumeType(t *testing.T) {
	createVolumeParameters := blockstorage.CreateVolumeTypeParameters{
		Name:       "displayname",
		ExtraSpecs: &blockstorage.VolumeTypeExtraSpecs{Capabilities: "gpu"},
	}

	createVolumeTypeTest(t, createVolumeParameters, `{"volume_type":{"name":"displayname","extra_specs":{"capabilities":"gpu"}}}`)
}

func TestCreateVolumeTypeWithOptionalExtraSpecsNotInPayload(t *testing.T) {
	createVolumeParameters := blockstorage.CreateVolumeTypeParameters{Name: "displayname"}

	createVolumeTypeTest(t, createVolumeParameters, `{"volume_type":{"name":"displayname"}}`)
}

func createVolumeTypeTest(t *testing.T, p blockstorage.CreateVolumeTypeParameters, expectedPayload string) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleVolumeTypeJson, "types", expectedPayload)
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	volume, err := service.CreateVolumeType(p)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, volumeTypeSample, volume)
}
