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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleFlavor = Flavor{
	ID:   101,
	Name: "standard.small",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}}}

var sampleFlavor2 = Flavor{
	ID:   102,
	Name: "standard.medium",
	Links: []Link{
		{"href_1", "rel_1"},
		{"href_2", "rel_2"}}}

func TestGetFlavors(t *testing.T) {

	mockResponseObject := FlavorsContainer{
		Flavors: []Flavor{sampleFlavor, sampleFlavor2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/flavors")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	flavors, err := service.Flavors()
	testUtil.IsNil(t, err)

	if len(flavors) != 2 {
		t.Error(errors.New("Error: Expected 2 flavors to be listed"))
	}
	testUtil.Equals(t, sampleFlavor, flavors[0])
	testUtil.Equals(t, sampleFlavor2, flavors[1])
}

func TestGetFlavor(t *testing.T) {

	flavorID := "101"
	mockResponseObject := FlavorContainer{Flavor: sampleFlavor}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, misc.Strcat("/flavors/", flavorID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	result, err := service.Flavor(flavorID)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, sampleFlavor, result)
}

func TestGetFlavorsInvalid(t *testing.T) {

	mockResponseObject := FlavorsContainer{
		Flavors: []Flavor{sampleFlavor, sampleFlavor2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, "/flavors")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Flavors()
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetFlavorWithInvalidPayload(t *testing.T) {

	flavorID := "101"
	mockResponseObject := FlavorContainer{Flavor: sampleFlavor}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, misc.Strcat("/flavors/", flavorID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Flavor(flavorID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetFlavorWithInvalidFlavorID(t *testing.T) {

	invalidFlavorID := "999"
	mockResponseObject := FlavorContainer{Flavor: sampleFlavor}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, http.StatusNotFound, mockResponseObject, misc.Strcat("/flavors/", invalidFlavorID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Flavor(invalidFlavorID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
	status, ok := err.(misc.HTTPStatus)
	if ok {
		testUtil.Equals(t, http.StatusNotFound, status.StatusCode)
	}
}

func TestParseFlavors(t *testing.T) {

	flavorsTestFilePath := "./testdata/database_flavor_test_flavors.json"

	flavorsTestFile, err := os.Open("./testdata/database_flavor_test_flavors.json")
	if err != nil {
		t.Error(fmt.Errorf("Failed to open file %s: '%s'", flavorsTestFilePath, err.Error()))
	}

	flavorsContainer := FlavorsContainer{}
	err = json.NewDecoder(flavorsTestFile).Decode(&flavorsContainer)
	defer flavorsTestFile.Close()
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", flavorsTestFilePath, err.Error()))
	}

	flavorTestFilePath := "./testdata/database_flavor_test_flavor.json"

	flavorTestFile, err := os.Open(flavorTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to open file %s: '%s'", flavorTestFilePath, err.Error()))
	}

	flavorContainer := FlavorContainer{}
	err = json.NewDecoder(flavorTestFile).Decode(&flavorContainer)
	defer flavorTestFile.Close()
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", flavorTestFilePath, err.Error()))
	}

	return
}
