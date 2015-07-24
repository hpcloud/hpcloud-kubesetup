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
	"errors"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var instanceID = "829144e8-0a03-40a5-8bc9-38756b111111"
var sampleDatabase1 = Database{Name: "test1"}
var sampleDatabase2 = Database{Name: "test2"}

func TestGetDatabases(t *testing.T) {

	mockResponseObject := DatabasesContainer{
		Databases: []Database{sampleDatabase1, sampleDatabase2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/databases")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	databases, err := service.Databases(instanceID)
	if err != nil {
		t.Error(err)
	}

	if len(databases) != 2 {
		t.Error(errors.New("Error: Expected 2 databases to be listed"))
	}
	testUtil.Equals(t, sampleDatabase1, databases[0])
	testUtil.Equals(t, sampleDatabase2, databases[1])
}

func TestGetDatabasesInvalid(t *testing.T) {

	mockResponseObject := DatabasesContainer{
		Databases: []Database{sampleDatabase1, sampleDatabase2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, "/databases")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Databases(instanceID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestDeleteDatabase(t *testing.T) {
	name := "user"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/instanceID/databases/"+name)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.DeleteDatabase("instanceID", name)
	testUtil.IsNil(t, err)
}

func TestCreateDatabase(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, "", "/instances/InstanceID/databases",
		`{"databases":[{"character_set":"utf32","collate":"latin","name":"username"}]}`)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.CreateDatabases("InstanceID", CreateDatabaseParameters{Name: "username", CharacterSet: "utf32", Collate: "latin"})
	testUtil.IsNil(t, err)
}
