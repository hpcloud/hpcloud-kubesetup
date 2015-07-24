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

package compute

import (
	"encoding/json"
	"errors"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleVolumeAttachment = VolumeAttachment{
	ID:       "id",
	Device:   "/dev/vdd",
	ServerID: "serverID",
	VolumeID: "volumeID"}

func TestGetVolumeAttachments(t *testing.T) {

	mockResponseObject := volumeAttachmentsResp{VolumeAttachments: []VolumeAttachment{sampleVolumeAttachment}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/servers/serverID/os-volume_attachments")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	attachedVolumes, err := service.VolumeAttachments("serverID")
	if err != nil {
		t.Fatal(err)
	}

	if len(attachedVolumes) != 1 {
		t.Fatal(errors.New("Error: Expected 1 volume to be attached"))
	}

	testUtil.Equals(t, sampleVolumeAttachment, attachedVolumes[0])

}

func TestGetVolumeAttachment(t *testing.T) {

	mockResponseObject := volumeAttachmentContainer{VolumeAttachment: sampleVolumeAttachment}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/servers/serverID/os-volume_attachments/id")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	attachedVolume, err := service.VolumeAttachment("serverID", "id")
	if err != nil {
		t.Fatal(err)
	}

	testUtil.Equals(t, sampleVolumeAttachment, attachedVolume)
}

func TestAttachVolume(t *testing.T) {
	mockResponseObject, _ := json.Marshal(volumeAttachmentContainer{VolumeAttachment: sampleVolumeAttachment})
	mockRequestObject, _ := json.Marshal(volumeAttachmentContainer{VolumeAttachment: VolumeAttachment{Device: "/dev/vdd", VolumeID: "volumeID"}})
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, string(mockResponseObject),
		"/servers/serverID/os-volume_attachments", string(mockRequestObject))
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	attachedVolume, err := service.AttachVolume("serverID", "volumeID", "/dev/vdd")
	if err != nil {
		t.Fatal(err)
	}

	testUtil.Equals(t, sampleVolumeAttachment, attachedVolume)
}

func TestDeleteVolumeAttachment(t *testing.T) {
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/servers/serverID/os-volume_attachments/Id")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	err := service.DeleteVolumeAttachment("serverID", "Id")
	if err != nil {
		t.Fatal(err)
	}
}
