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

import "git.openstack.org/stackforge/golang-client.git/misc"

// VolumeAttachment contains properties a volume
// that can be attached to a server and can
// be in a specified pool.
type VolumeAttachment struct {
	ID       string `json:"id,omitempty"`
	Device   string `json:"device"`
	ServerID string `json:"serverId,omitempty"`
	VolumeID string `json:"volumeId"`
}

//VolumeAttachments lists the volume attachments for a specified server.
func (computeService Service) VolumeAttachments(serverID string) ([]VolumeAttachment, error) {
	var r = volumeAttachmentsResp{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/os-volume_attachments")
	if err != nil {
		return r.VolumeAttachments, err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.VolumeAttachments, err

}

//VolumeAttachment shows details for the specified volume attachment.
func (computeService Service) VolumeAttachment(serverID string, attachmentID string) (VolumeAttachment, error) {
	var r = volumeAttachmentContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/os-volume_attachments/", attachmentID)
	if err != nil {
		return r.VolumeAttachment, err
	}
	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.VolumeAttachment, err
}

//AttachVolume attaches a volume to the specified server.
func (computeService Service) AttachVolume(serverID string, volumeID string, device string) (VolumeAttachment, error) {
	attachVolumeRequest := VolumeAttachment{Device: device, VolumeID: volumeID}
	input := volumeAttachmentContainer{VolumeAttachment: attachVolumeRequest}
	output := volumeAttachmentContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/os-volume_attachments")

	if err != nil {
		return output.VolumeAttachment, err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, input, &output)

	return output.VolumeAttachment, err
}

//DeleteVolumeAttachment deletes the specified volume attachment from a specified server.
func (computeService Service) DeleteVolumeAttachment(serverID string, attachmentID string) (err error) {
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/os-volume_attachments/", attachmentID)

	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)

}

type volumeAttachmentsResp struct {
	VolumeAttachments []VolumeAttachment `json:"volumeAttachments,omitempty"`
}

type volumeAttachmentContainer struct {
	VolumeAttachment VolumeAttachment `json:"volumeAttachment,omitempty"`
}
