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

// VolumeType has properties for a volume type in
// block storage in openstack.
type VolumeType struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	ExtraSpecs VolumeTypeExtraSpecs `json:"extra_specs"`
}

// VolumeTypeExtraSpecs has properties for a volumetype
// that has extra specs like gpu
type VolumeTypeExtraSpecs struct {
	Capabilities string `json:"capabilities,omitempty"`
}

// CreateVolumeTypeParameters are the properties for creating a new volume
type CreateVolumeTypeParameters struct {
	Name       string                `json:"name"`
	ExtraSpecs *VolumeTypeExtraSpecs `json:"extra_specs,omitempty"` // optional property so allowing this to be nil so as not to be in the payload.
}

// VolumeTypes will issue a GET request to retrieve the volume types.
func (blockStorageService Service) VolumeTypes() ([]VolumeType, error) {
	c := map[string][]VolumeType{"json:volumes": []VolumeType{}}
	url, err := blockStorageService.buildRequestURL("/types")
	if err != nil {
		return c["volume_types"], err
	}

	err = misc.GetJSON(url, blockStorageService.authenticator, &c)
	return c["volume_types"], err
}

// VolumeType will issue a GET request to retrieve the volumetype by id.
func (blockStorageService Service) VolumeType(id string) (VolumeType, error) {
	c := volumeTypeContainer{VolumeType: VolumeType{}}
	url, err := blockStorageService.buildRequestURL("/types/", id)
	if err != nil {
		return c.VolumeType, err
	}

	err = misc.GetJSON(url, blockStorageService.authenticator, &c)
	return c.VolumeType, err
}

// CreateVolumeType will send a POST request to create a new volumetype with the specified parameters.
func (blockStorageService Service) CreateVolumeType(parameters CreateVolumeTypeParameters) (VolumeType, error) {
	cIn := volumeTypeCreateParametersContainer{VolumeType: parameters}
	cOut := volumeTypeContainer{VolumeType: VolumeType{}}
	reqURL, err := blockStorageService.buildRequestURL("/types")
	if err != nil {
		return cOut.VolumeType, err
	}

	err = misc.PostJSON(reqURL, blockStorageService.authenticator, cIn, &cOut)
	return cOut.VolumeType, err
}

// DeleteVolumeType will delete the volume type by id.
func (blockStorageService Service) DeleteVolumeType(id string) (err error) {
	reqURL, err := blockStorageService.buildRequestURL("/types/", id)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, blockStorageService.authenticator)
}

type volumeTypeContainer struct {
	VolumeType VolumeType `json:"volume_type"`
}

type volumeTypeCreateParametersContainer struct {
	VolumeType CreateVolumeTypeParameters `json:"volume_type"`
}
