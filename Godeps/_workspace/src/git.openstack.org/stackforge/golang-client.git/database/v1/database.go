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

// Database contains information of OpenStack trove database
type Database struct {
	Name string `json:"name"`
}

// DatabasesContainer contains information for databases.
type DatabasesContainer struct {
	Databases []Database `json:"databases"`
}

// DBContainer contains information for a database.
type DBContainer struct {
	Database Database `json:"database"`
}

// CreateDatabaseParameters has properties for used when creating a database instance.
type CreateDatabaseParameters struct {
	CharacterSet string `json:"character_set,omitempty"`
	Collate      string `json:"collate,omitempty"`
	Name         string `json:"name"`
}

// Databases will issue a GET request to retrieve all databases.
func (databaseService Service) Databases(instanceID string) ([]Database, error) {

	var container = DatabasesContainer{}
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/databases")
	if err != nil {
		return container.Databases, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)

	return container.Databases, err
}

// CreateDatabases will issue a POST request to create the specified databases.
func (databaseService Service) CreateDatabases(instanceID string, params ...CreateDatabaseParameters) error {

	var inContainer = createDatabasesParametersContainer{params}
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/databases")
	if err != nil {
		return err
	}

	return misc.PostJSON(reqURL, databaseService.authenticator, &inContainer, nil)
}

// DeleteDatabase will issue a delete query to delete the database.
func (databaseService Service) DeleteDatabase(instanceID string, databaseName string) (err error) {
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/databases/", databaseName)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, databaseService.authenticator)
}

type createDatabasesParametersContainer struct {
	Databases []CreateDatabaseParameters `json:"databases"`
}
