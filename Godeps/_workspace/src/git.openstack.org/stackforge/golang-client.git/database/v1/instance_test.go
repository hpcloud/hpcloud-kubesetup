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

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var tokn = "test-token-1"
var createdTime, _ = misc.NewDateTimeFromString(`"2014-12-05T10:30:30"`)
var updatedTime, _ = misc.NewDateTimeFromString(`"2014-12-08T11:30:30"`)

var sampleInstance1 = Instance{
	Status: "ACTIVE",
	Name:   "db1",
	Links: []compute.Link{
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/instances/829144e8-0a03-40a5-8bc9-38756b111111",
			Rel: "self"},
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/instances/829144e8-0a03-40a5-8bc9-38756b111111",
			Rel: "bookmark"}},
	IP: []string{"10.9.194.40", "15.125.35.100"},
	ID: "829144e8-0a03-40a5-8bc9-38756b111111",
	Flavor: compute.Flavor{
		ID:   "1001",
		Name: "db.xsmall",
		Links: []compute.Link{
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/flavors/1001",
				Rel: "self"},
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/flavors/1001",
				Rel: "bookmark"}}},
	Datastore: Datastore{Version: "5.1", Type: "mysql"}}

var sampleInstance2 = Instance{
	Status: "ACTIVE",
	Name:   "db2",
	Links: []compute.Link{
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/instances/829144e8-0a03-40a5-8bc9-38756b222222",
			Rel: "self"},
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/instances/829144e8-0a03-40a5-8bc9-38756b222222",
			Rel: "bookmark"}},
	IP: []string{"10.9.194.41", "15.125.35.200"},
	ID: "829144e8-0a03-40a5-8bc9-38756b222222",
	Flavor: compute.Flavor{
		ID:   "1002",
		Name: "db.small",
		Links: []compute.Link{
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/flavors/1002",
				Rel: "self"},
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/flavors/1002",
				Rel: "bookmark"}}},
	Datastore: Datastore{Version: "5.5", Type: "mysql"}}

var sampleInstanceDetail = InstanceDetail{
	Status:  "ACTIVE",
	Updated: updatedTime,
	Name:    "db1",
	Links: []compute.Link{
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/instances/829144e8-0a03-40a5-8bc9-38756b111111",
			Rel: "self"},
		{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/instances/829144e8-0a03-40a5-8bc9-38756b111111",
			Rel: "bookmark"}},
	Created:      createdTime,
	IP:           []string{"10.9.194.40", "15.125.35.100"},
	LocalStorage: LocalStorage{Size: 1000, Used: 0.55},
	ID:           "829144e8-0a03-40a5-8bc9-38756b111111",
	Flavor: compute.Flavor{
		ID:   "1001",
		Name: "db.xsmall",
		Links: []compute.Link{
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/v1.0/10857480012345/flavors/1001",
				Rel: "self"},
			{HRef: "https://region-a.geo-1.database.hpcloudsvc.com/flavors/1001",
				Rel: "bookmark"}}},
	Datastore: Datastore{Version: "5.1", Type: "mysql"}}

func TestGetInstances(t *testing.T) {

	mockResponseObject := InstancesContainer{
		Instances: []Instance{sampleInstance1, sampleInstance2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/instances")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	instances, err := service.Instances()
	if err != nil {
		t.Error(err)
	}

	if len(instances) != 2 {
		t.Error(errors.New("Error: Expected 2 instances to be listed"))
	}
	testUtil.Equals(t, sampleInstance1, instances[0])
	testUtil.Equals(t, sampleInstance2, instances[1])
}

func TestGetInstance(t *testing.T) {

	instanceID := "829144e8-0a03-40a5-8bc9-38756b111111"
	mockResponseObject := InstanceDetailContainer{Instance: sampleInstanceDetail}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, misc.Strcat("/instances/", instanceID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	result, err := service.Instance(instanceID)
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleInstanceDetail, result)
}

func TestDeleteInstance(t *testing.T) {
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "instances/instanceID")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.DeleteInstance("instanceID")
	testUtil.IsNil(t, err)
}

var sampleInstanceDetailBytes, _ = json.Marshal(sampleInstanceDetail)
var sampleInstanceDetailResponse = `{ "instance":` + string(sampleInstanceDetailBytes) + `}`

func TestCreateInstance(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleInstanceDetailResponse, "/instances",
		`{"instance":{"databases":[{"name":"newDb"}],"flavorRef":"1001","name":"Instance1","users":[{"name":"User1","password":"ag92340gv"}],"nics":[{"net-id":"net-id-1234"}]}}`)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	user := UserParameter{Name: "User1", Password: "ag92340gv"}
	db := CreateDatabaseParameters{
		Name: "newDb",
	}

	instance, err := service.CreateInstance(CreateInstanceParameters{
		Databases: []CreateDatabaseParameters{db},
		FlavorRef: "1001",
		Name:      "Instance1",
		Users:     []UserParameter{user},
		NetworkInterfaces: []NetworkInterface{
			NetworkInterface{NetID: "net-id-1234"},
		},
	})

	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleInstanceDetail, instance)
}

func TestRestoreInstance(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleInstanceDetailResponse, "/instances",
		`{"instance":{"flavorRef":"1001","restorePoint":{"backupRef":"438hg234"},"name":"Instance1"}}`)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	instance, err := service.RestoreBackupInstance(RestoreBackupParameters{
		Name:         "Instance1",
		RestorePoint: RestorePoint{BackupRef: "438hg234"},
		FlavorRef:    "1001"})

	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleInstanceDetail, instance)
}

func TestEnableRoot(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, `{"user": {"password": "3t42g", "name": "root"}}`, "/instances/InstanceID/root", ``)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	val, err := service.EnableRoot("InstanceID")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "root", val.Name)
	testUtil.Equals(t, "3t42g", val.Password)
}

func TestGetInstancesInvalid(t *testing.T) {

	mockResponseObject := InstancesContainer{
		Instances: []Instance{sampleInstance1, sampleInstance2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, "/instances")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Instances()
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetInstanceWithInvalidPayload(t *testing.T) {

	instanceID := "829144e8-0a03-40a5-8bc9-38756b111111"
	mockResponseObject := InstanceContainer{Instance: sampleInstance1}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, misc.Strcat("/instances/", instanceID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Instance(instanceID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestGetInstanceWithInvalidInstanceID(t *testing.T) {

	invalidInstanceID := "999"
	mockResponseObject := InstanceDetailContainer{Instance: sampleInstanceDetail}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, http.StatusNotFound, mockResponseObject, misc.Strcat("/instances/", invalidInstanceID))
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Instance(invalidInstanceID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
	status, ok := err.(misc.HTTPStatus)
	if ok {
		testUtil.Equals(t, http.StatusNotFound, status.StatusCode)
	}
}

func TestParseInstances(t *testing.T) {

	instancesTestFilePath := "./testdata/instance_test_instances.json"
	instancesTestFile, err := os.Open(instancesTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to open file %s: '%s'", instancesTestFilePath, err.Error()))
	}

	instancesContainer := InstancesContainer{}
	err = json.NewDecoder(instancesTestFile).Decode(&instancesContainer)
	defer instancesTestFile.Close()
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", instancesTestFilePath, err.Error()))
	}

	instanceTestFilePath := "./testdata/instance_test_instance.json"
	instanceTestFile, err := os.Open(instanceTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to open file %s: '%s'", instanceTestFilePath, err.Error()))
	}

	instanceContainer := InstanceDetailContainer{}
	err = json.NewDecoder(instanceTestFile).Decode(&instanceContainer)
	defer instanceTestFile.Close()
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", instanceTestFilePath, err.Error()))
	}

	return
}

func CreateDatabaseService(url string) Service {
	return NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: url})
}
