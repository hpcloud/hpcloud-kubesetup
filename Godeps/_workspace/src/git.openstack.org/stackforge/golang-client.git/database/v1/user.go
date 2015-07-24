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

// User contains information of OpenStack trove database user
type User struct {
	Databases []Database `json:"database"`
	Name      string     `json:"name"`
}

// UsersContainer contains information for users.
type UsersContainer struct {
	Users []User `json:"users"`
}

// UserContainer contains information for a user.
type UserContainer struct {
	User User `json:"user"`
}

// UserParameter has properties for a database user.
type UserParameter struct {
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	Databases []Database `json:"databases,omitempty"`
}

// Users will issue a GET request to retrieve information of all users.
func (databaseService Service) Users(instanceID string) ([]User, error) {

	var container = UsersContainer{}
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/users")
	if err != nil {
		return container.Users, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)
	return container.Users, err
}

// CreateUser will issue a POST request to create users for the instance.
func (databaseService Service) CreateUser(instanceID string, params ...UserParameter) error {

	var inContainer = createUsersParametersContainer{params}
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/users")
	if err != nil {
		return err
	}

	return misc.PostJSON(reqURL, databaseService.authenticator, &inContainer, nil)
}

// DeleteUser will issue a delete query to delete the database.
func (databaseService Service) DeleteUser(instanceID string, userName string) (err error) {
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/users/", userName)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, databaseService.authenticator)
}

type createUsersParametersContainer struct {
	Users []UserParameter `json:"users"`
}

type userParameterContainer struct {
	User UserParameter `json:"user"`
}
