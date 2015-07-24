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

var sampleUser1 = User{Name: "user1", Databases: []Database{Database{Name: "test1"}, Database{Name: "test2"}}}
var sampleUser2 = User{Name: "user2", Databases: []Database{Database{Name: "test2"}}}

func TestGetUsers(t *testing.T) {

	mockResponseObject := UsersContainer{
		Users: []User{sampleUser1, sampleUser2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/users")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	users, err := service.Users(instanceID)
	if err != nil {
		t.Error(err)
	}

	if len(users) != 2 {
		t.Error(errors.New("Error: Expected 2 users to be listed"))
	}
	testUtil.Equals(t, sampleUser1, users[0])
	testUtil.Equals(t, sampleUser2, users[1])
}

func TestGetUsersInvalid(t *testing.T) {

	mockResponseObject := UsersContainer{
		Users: []User{sampleUser1, sampleUser2}}

	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, tokn, testUtil.InvalidJSONPayload, mockResponseObject, "/users")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	_, err := service.Users(instanceID)
	if err == nil {
		t.Error(errors.New("Error: Expected error was not returned."))
	}
}

func TestDeleteUser(t *testing.T) {
	name := "user"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/instanceID/users/"+name)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.DeleteUser("instanceID", name)
	testUtil.IsNil(t, err)
}

func TestCreateUser(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, "", "/instances/InstanceID/users",
		`{"users":[{"name":"username","password":"32thbw3"}]}`)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.CreateUser("InstanceID", UserParameter{Name: "username", Password: "32thbw3"})
	testUtil.IsNil(t, err)
}
