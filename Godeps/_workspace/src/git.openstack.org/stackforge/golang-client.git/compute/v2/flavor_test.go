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
	"errors"
	"net/http"

	"testing"

	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleFlavor = Flavor{
	ID:   "101",
	Name: "standard.small",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}}}

var sampleFlavor2 = Flavor{
	ID:   "102",
	Name: "standard.medium",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}}}

var sampleFlavorDetail = FlavorDetail{
	ID:   "101",
	Name: "standard.small",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}},
	RAM:        misc.Int64Wrapper{Int64: 512, Valid: true},
	VCPUs:      misc.Int64Wrapper{Int64: 1, Valid: true},
	Disk:       misc.Int64Wrapper{Int64: 10, Valid: true},
	Ephemeral:  misc.Int64Wrapper{Int64: 5, Valid: true},
	RXTXFactor: misc.Float64Wrapper{Float64: 1.3, Valid: true}}

var sampleFlavorDetail2 = FlavorDetail{
	ID:   "102",
	Name: "standard.medium",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}},
	RAM:        misc.Int64Wrapper{Int64: 1024, Valid: true},
	VCPUs:      misc.Int64Wrapper{Int64: 2, Valid: true},
	Disk:       misc.Int64Wrapper{Int64: 10, Valid: true},
	Ephemeral:  misc.Int64Wrapper{Int64: 5, Valid: true},
	RXTXFactor: misc.Float64Wrapper{Float64: 1.5, Valid: true}}

func TestGetFlavors(t *testing.T) {
	mockResponseObject := FlavorsContainer{
		Flavors: []Flavor{sampleFlavor, sampleFlavor2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/flavors?limit=21")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	flavors, err := service.Flavors()
	testUtil.IsNil(t, err)

	if len(flavors) != 2 {
		t.Error(errors.New("Error: Expected 2 flavors to be listed"))
	}
	testUtil.Equals(t, sampleFlavor, flavors[0])
	testUtil.Equals(t, sampleFlavor2, flavors[1])
}

func TestGetFlavorsDetail(t *testing.T) {
	mockResponseObject := FlavorsDetailContainer{
		FlavorsDetail: []FlavorDetail{sampleFlavorDetail, sampleFlavorDetail2}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/flavors/detail?limit=21")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	flavorsDetail, err := service.FlavorsDetail()
	testUtil.IsNil(t, err)

	if len(flavorsDetail) != 2 {
		t.Error(errors.New("Error: Expected 2 flavors detail to be listed"))
	}
	testUtil.Equals(t, sampleFlavorDetail, flavorsDetail[0])
	testUtil.Equals(t, sampleFlavorDetail2, flavorsDetail[1])
}

func TestGetFlavorDetail(t *testing.T) {
	flavorID := "101"
	mockResponseObject := FlavorDetailContainer{FlavorDetail: sampleFlavorDetail}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, misc.Strcat("/flavors/", flavorID))
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	result, err := service.FlavorDetail(flavorID)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, sampleFlavorDetail, result)
}

func TestGetFlavorsInvalid(t *testing.T) {
	mockResponseObject := FlavorsContainer{
		Flavors: []Flavor{sampleFlavor, sampleFlavor2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, "/flavors?limit=21")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	_, err := service.Flavors()
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetFlavorWithInvalidPayload(t *testing.T) {
	flavorID := "101"
	mockResponseObject := FlavorDetailContainer{FlavorDetail: sampleFlavorDetail}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, misc.Strcat("/flavors/", flavorID))
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	_, err := service.FlavorDetail(flavorID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetFlavorWithInvalidFlavorID(t *testing.T) {
	invalidFlavorID := "999"
	mockResponseObject := FlavorContainer{Flavor: sampleFlavor}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, http.StatusNotFound, mockResponseObject, misc.Strcat("/flavors/", invalidFlavorID))
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	_, err := service.FlavorDetail(invalidFlavorID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
	status, ok := err.(misc.HTTPStatus)
	if ok {
		testUtil.Equals(t, http.StatusNotFound, status.StatusCode)
	}
}
