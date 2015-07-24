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

package database

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Flavor contains information of OpenStack trove flavor
// TODO: we need to consider consolidating the two Flavor structs (one in compute, and one in database)
//       since these two structs are only slightly different: ID in one is string while ID is int32 in the other.
type Flavor struct {
	ID    int32  `json:"id"`
	Links []Link `json:"links"`
	Name  string `json:"name,omitempty"`
}

// Link contains information of a link
type Link struct {
	HRef string `json:"href"`
	Rel  string `json:"rel"`
}

// FlavorsContainer contains information for flavors.
type FlavorsContainer struct {
	Flavors []Flavor `json:"flavors"`
}

// FlavorContainer contains information for a flavor.
type FlavorContainer struct {
	Flavor Flavor `json:"flavor"`
}

// Flavors will issue a GET request to retrieve all flavors.
func (databaseService Service) Flavors() ([]Flavor, error) {

	var container = FlavorsContainer{}
	reqURL, err := databaseService.buildRequestURL("/flavors")
	if err != nil {
		return container.Flavors, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)

	return container.Flavors, err
}

// Flavor will issue a GET request to retrieve the specified flavor.
func (databaseService Service) Flavor(flavorID string) (Flavor, error) {

	var container = FlavorContainer{}
	reqURL, err := databaseService.buildRequestURL("/flavors/", flavorID)
	if err != nil {
		return container.Flavor, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)
	return container.Flavor, err
}
