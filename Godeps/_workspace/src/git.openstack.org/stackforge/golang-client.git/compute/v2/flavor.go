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
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Flavor contains information of OpenStack nova from the types of servers.
type Flavor struct {
	ID    string `json:"id"`
	Links []Link `json:"links"`
	Name  string `json:"name,omitempty"`
}

// FlavorDetail contains information of OpenStack nova flavor detail
type FlavorDetail struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Links      []Link              `json:"links"`
	RAM        misc.Int64Wrapper   `json:"ram"`
	VCPUs      misc.Int64Wrapper   `json:"vcpus"`
	Swap       misc.Int64Wrapper   `json:"swap"`
	Disk       misc.Int64Wrapper   `json:"disk"`
	Ephemeral  misc.Int64Wrapper   `json:"OS-FLV-EXT-DATA:ephemeral"`
	RXTXFactor misc.Float64Wrapper `json:"rxtx_factor"`
}

// FlavorsContainer contains information for flavors.
type FlavorsContainer struct {
	Flavors []Flavor `json:"flavors"`
}

// FlavorContainer contains information for a flavor.
type FlavorContainer struct {
	Flavor Flavor `json:"flavor"`
}

// FlavorsDetailContainer contains information for flavors detail.
type FlavorsDetailContainer struct {
	FlavorsDetail []FlavorDetail `json:"flavors"`
}

// FlavorDetailContainer contains information for a flavor detail.
type FlavorDetailContainer struct {
	FlavorDetail FlavorDetail `json:"flavor"`
}

// Flavors will issue a GET request to retrieve all flavors.
func (computeService Service) Flavors() ([]Flavor, error) {
	result := []Flavor{}
	marker := ""
	more := true

	var err error
	for more {
		container := FlavorsContainer{}
		url, err := computeService.buildPaginatedQueryURL(QueryParameters{Limit: defaultLimit, Marker: marker}, "/flavors")
		if err != nil {
			return container.Flavors, err
		}

		err = misc.GetJSON(url.String(), computeService.authenticator, &container)
		if err != nil {
			return nil, err
		}

		if len(container.Flavors) < defaultLimit {
			more = false
		}

		for _, flavor := range container.Flavors {
			result = append(result, flavor)
			marker = flavor.ID
		}
	}

	return result, err
}

// FlavorsDetail will issue a GET request to retrieve all flavors detail.
func (computeService Service) FlavorsDetail() ([]FlavorDetail, error) {
	result := []FlavorDetail{}
	marker := ""
	more := true

	var err error
	for more {
		container := FlavorsDetailContainer{}
		url, err := computeService.buildPaginatedQueryURL(QueryParameters{Limit: defaultLimit, Marker: marker}, "/flavors/detail")
		if err != nil {
			return container.FlavorsDetail, err
		}

		err = misc.GetJSON(url.String(), computeService.authenticator, &container)
		if err != nil {
			return nil, err
		}

		if len(container.FlavorsDetail) < defaultLimit {
			more = false
		}

		for _, flavorDetail := range container.FlavorsDetail {
			result = append(result, flavorDetail)
			marker = flavorDetail.ID
		}
	}

	return result, err
}

// FlavorDetail will issue a GET request to retrieve the specified flavor detail.
func (computeService Service) FlavorDetail(flavorID string) (FlavorDetail, error) {
	var container = FlavorDetailContainer{}
	url, err := computeService.buildRequestURL("/flavors/", flavorID)
	if err != nil {
		return container.FlavorDetail, err
	}

	err = misc.GetJSON(url, computeService.authenticator, &container)
	return container.FlavorDetail, err
}
