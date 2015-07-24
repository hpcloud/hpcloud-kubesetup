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
	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var createdAt, _ = misc.NewDateTimeFromString(`"2014-09-29T15:33:37"`)
var volumeSample = blockstorage.Volume{
	Status:      "ACTIVE",
	DisplayName: "Display",
	Attachments: []blockstorage.VolumeAttachment{
		blockstorage.VolumeAttachment{
			AttachmentID: "attachment-id-1234",
			Device:       "/dev/vdc",
			ServerID:     "server-id-1234",
			VolumeID:     "ID",
		},
	},
	Az:          "az1",
	Bootable:    true,
	CreatedAt:   createdAt,
	Desciption:  "Description",
	VolumeType:  "VolType",
	SnapshotID:  "SnapshotID",
	SourceVolID: "sourceVolid",
	Metadata:    map[string]string{"metaItem": "foo"},
	ID:          "ID",
	Size:        13241,
	ImageRef:    "f76101bc-c442-4619-a09a-5296748e12f8",
}

var volumeJSON, _ = json.Marshal(volumeSample)

var sampleVolumesJSON = fmt.Sprintf(`{"volumes" : [ %s ] }`, string(volumeJSON))

var sampleVolumeJSON = fmt.Sprintf(`{"volume" : %s }`, string(volumeJSON))

var tokn = "2926072626"

func TestGetVolumes(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyStatusAndURL(t, tokn, 200, sampleVolumesJSON, "/volumes")
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	volumes, err := service.Volumes()
	testUtil.IsNil(t, err)
	testUtil.Assert(t, len(volumes) == 1, "Expected 1 volume")
	testUtil.Equals(t, volumeSample, volumes[0])
}

func TestGetVolume(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyStatusAndURL(t, tokn, 200, sampleVolumeJSON, "/volumes/id")
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	volume, err := service.Volume("id")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, volumeSample, volume)
}

func TestDeleteVolume(t *testing.T) {
	name := "id"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "volumes/"+name)
	defer apiServer.Close()

	service := CreateVolumeService(apiServer.URL)
	err := service.DeleteVolume(name)
	testUtil.IsNil(t, err)
}

func TestCreateVolume(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleVolumeJSON, "volumes",
		`{"volume":{"display_name":"displayname","availability_zone":"az","display_description":"Description","volume_type":"volType","size":512315}}`)
	defer apiServer.Close()

	createVolumeParameters := blockstorage.CreateVolumeParameters{
		DisplayName: "displayname",
		Az:          "az",
		Desciption:  "Description",
		VolumeType:  "volType",
		Metadata:    map[string]string{},
		Size:        512315,
	}

	service := CreateVolumeService(apiServer.URL)
	volume, err := service.CreateVolume(createVolumeParameters)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, volumeSample, volume)
}

func CreateVolumeService(url string) blockstorage.Service {
	return blockstorage.NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: url})
}
