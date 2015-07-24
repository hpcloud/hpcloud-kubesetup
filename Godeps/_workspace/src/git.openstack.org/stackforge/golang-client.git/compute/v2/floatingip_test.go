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
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var samplefloatingIP = compute.FloatingIP{
	ID:         "id",
	IP:         "14.14.154.15",
	InstanceID: "instanceid",
	FixedIP:    "fixedip",
	Pool:       "pool"}

func TestGetFloatingIPs(t *testing.T) {
	mockResponseObject := floatingIPsContainer{FloatingIPs: []compute.FloatingIP{samplefloatingIP}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "os-floating-ips")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	floatingIPs, err := service.FloatingIPs()
	if err != nil {
		t.Error(err)
	}

	if len(floatingIPs) != 1 {
		t.Error(errors.New("Error: Expected 1 floating ip to be listed"))
	}
	testUtil.Equals(t, samplefloatingIP, floatingIPs[0])
}

func TestGetFloatingIP(t *testing.T) {
	mockResponseObject := floatingIPContainer{FloatingIP: samplefloatingIP}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "os-floating-ips/id")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	keypair, err := service.FloatingIP("id")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, samplefloatingIP, keypair)
}

func TestDeleteFloatingIP(t *testing.T) {
	name := "keypairName"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "os-floating-ips/"+name)
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	err := service.DeleteFloatingIP(name)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateFloatingIP(t *testing.T) {
	mockResponse, _ := json.Marshal(floatingIPContainer{FloatingIP: samplefloatingIP})
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, string(mockResponse), "os-floating-ips",
		`{"pool":"poolName"}`)
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	floatingIP, err := service.CreateFloatingIP("poolName")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, samplefloatingIP, floatingIP)
}

type floatingIPsContainer struct {
	FloatingIPs []compute.FloatingIP `json:"floating_ips"`
}

type floatingIPContainer struct {
	FloatingIP compute.FloatingIP `json:"floating_ip"`
}
