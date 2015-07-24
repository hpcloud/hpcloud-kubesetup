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

package blockstorage

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Volume has properties for a volume in
// block storage in openstack.
type Volume struct {
	Status      string               `json:"status"`
	DisplayName string               `json:"display_name"`
	Attachments []VolumeAttachment   `json:"attachments"`
	Az          string               `json:"availability_zone"`
	Bootable    bool                 `json:"bootable,string"`
	Encrypted   bool                 `json:"encrypted,omitempty"` // property not seen in public Helion
	CreatedAt   misc.RFC8601DateTime `json:"created_at"`
	Desciption  string               `json:"display_description"`
	VolumeType  string               `json:"volume_type"`
	SnapshotID  string               `json:"snapshot_id"`
	SourceVolID string               `json:"source_volid"`
	Metadata    map[string]string    `json:"metadata,omitempty"`
	ID          string               `json:"id"`
	Size        int                  `json:"size"`
	ImageRef    string               `json:"imageRef,omitempty"`
}

// VolumeAttachment describes what the volume is attached to
type VolumeAttachment struct {
	AttachmentID string `json:"id"`
	Device       string `json:"device"`
	ServerID     string `json:"server_id"`
	VolumeID     string `json:"volume_id"`
}

// CreateVolumeParameters are the properties for creating a new volume
type CreateVolumeParameters struct {
	DisplayName string            `json:"display_name"`
	Az          string            `json:"availability_zone,omitempty"`
	Desciption  string            `json:"display_description"`
	VolumeType  string            `json:"volume_type"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Size        int               `json:"size"`
	ImageRef    string            `json:"imageRef,omitempty"` //this is not documented but it works
}

// Volumes will issue a GET request to retrieve the volumes.
func (blockStorageService Service) Volumes() ([]Volume, error) {
	c := map[string][]Volume{"json:volumes": []Volume{}}
	url, err := blockStorageService.buildRequestURL("/volumes")
	if err != nil {
		return c["volumes"], err
	}

	err = misc.GetJSON(url, blockStorageService.authenticator, &c)
	return c["volumes"], err
}

// Volume will issue a GET request to retrieve the volume by id.
func (blockStorageService Service) Volume(id string) (Volume, error) {
	c := volumeContainer{Volume: Volume{}}
	url, err := blockStorageService.buildRequestURL("/volumes/", id)
	if err != nil {
		return c.Volume, err
	}

	err = misc.GetJSON(url, blockStorageService.authenticator, &c)
	return c.Volume, err
}

// CreateVolume will send a POST request to create a new volume with the specified parameters.
func (blockStorageService Service) CreateVolume(parameters CreateVolumeParameters) (Volume, error) {
	cIn := volumeCreateParametersContainer{Volume: parameters}
	cOut := volumeContainer{Volume: Volume{}}
	reqURL, err := blockStorageService.buildRequestURL("/volumes")
	if err != nil {
		return cOut.Volume, err
	}

	err = misc.PostJSON(reqURL, blockStorageService.authenticator, cIn, &cOut)
	return cOut.Volume, err
}

// DeleteVolume will delete the volume by id.
func (blockStorageService Service) DeleteVolume(id string) (err error) {
	reqURL, err := blockStorageService.buildRequestURL("/volumes/", id)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, blockStorageService.authenticator)
}

type volumeContainer struct {
	Volume Volume `json:"volume"`
}

type volumeCreateParametersContainer struct {
	Volume CreateVolumeParameters `json:"volume"`
}
