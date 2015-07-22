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

package compute_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestPrivateCloudAvailabilityZones(t *testing.T) {

	// Read private cloud response content from file
	availabilityZonesTestFilePath := "./testdata/availabilityzone_test_azs.json"
	availabilityZonesTestFileContent, err := ioutil.ReadFile(availabilityZonesTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to read JSON file %s: '%s'", availabilityZonesTestFilePath, err.Error()))
	}

	// Decode the content
	sampleAzs := compute.AvailabilityZonesContainer{}
	err = json.Unmarshal(availabilityZonesTestFileContent, &sampleAzs)
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", availabilityZonesTestFilePath, err.Error()))
	}

	// Test the SDK API computeService.AvailabilityZones()
	anon := func(computeService *compute.Service) {

		availabilityZones, err := computeService.AvailabilityZones()
		if err != nil {
			t.Error(err)
		}

		if len(availabilityZones) != len(sampleAzs.AvailabilityZones) {
			t.Error(errors.New("Incorrect number of availabilityZones found"))
		}

		// Verify returned availability zones match original sample availability zones
		testUtil.Equals(t, sampleAzs.AvailabilityZones, availabilityZones)
	}

	testComputeServiceAction(t, "os-availability-zone", string(availabilityZonesTestFileContent), anon)
}

func testComputeServiceAction(t *testing.T, uriEndsWith string, testData string, computeServiceAction func(*compute.Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	computeService := CreateComputeService(apiServer.URL)
	computeServiceAction(&computeService)
}
