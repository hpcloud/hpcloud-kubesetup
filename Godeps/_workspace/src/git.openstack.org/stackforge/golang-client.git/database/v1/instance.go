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
	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Database instance status.
const (
	BuildStatus           = "BUILD"
	RebootStatus          = "REBOOT"
	ActiveStatus          = "ACTIVE"
	FailedStatus          = "FAILED"
	BackupStatus          = "BACKUP"
	BlockedStatus         = "BLOCKED"
	ResizeStatus          = "RESIZE"
	RestartRequiredStatus = "RESTART_REQUIRED"
	ShutdownStatus        = "SHUTDOWN"
	ErrorStatus           = "ERROR"
)

// Service is a client service that can make
// requests against a OpenStack database v1 service.
type Service struct {
	authenticator common.Authenticator
}

// NewService a new Database service client
func NewService(authenticator common.Authenticator) Service {
	return Service{authenticator: authenticator}
}

func (databaseService Service) buildRequestURL(suffixes ...string) (string, error) {
	serviceURL, err := databaseService.authenticator.GetServiceURL("database", "1")
	if err != nil {
		return "", err
	}
	urlPaths := append([]string{serviceURL}, suffixes...)
	return misc.Strcat(urlPaths...), nil
}

// Instance contains information of an OpenStack trove database instance
type Instance struct {
	Status    string         `json:"status"`
	Name      string         `json:"name"`
	Links     []compute.Link `json:"links"`
	IP        []string       `json:"ip"`
	ID        string         `json:"id"`
	Flavor    compute.Flavor `json:"flavor"`
	Datastore Datastore      `json:"datastore"`
}

// InstanceDetail contains detailed information of an OpenStack trove database instance
type InstanceDetail struct {
	Status       string               `json:"status"`
	Updated      misc.RFC8601DateTime `json:"updated"`
	Name         string               `json:"name"`
	Links        []compute.Link       `json:"links"`
	Created      misc.RFC8601DateTime `json:"created"`
	IP           []string             `json:"ip"`
	LocalStorage LocalStorage         `json:"local_storage"`
	ID           string               `json:"id"`
	Flavor       compute.Flavor       `json:"flavor"`
	Datastore    Datastore            `json:"datastore"`
}

// LocalStorage contains information of OpenStack trove local storage
type LocalStorage struct {
	Size uint    `json:"size,omitempty"`
	Used float64 `json:"used"`
}

// Datastore contains information of OpenStack trove datastore
type Datastore struct {
	Version   string `json:"version"`
	Type      string `json:"type"`
	VersionID string `json:"version_id"`
}

// InstancesContainer is a container of information of database instances.
type InstancesContainer struct {
	Instances []Instance `json:"instances"`
}

// InstanceContainer is a container of information of a database instance.
type InstanceContainer struct {
	Instance Instance `json:"instance"`
}

// InstanceDetailContainer is a container of detailed information for a database instance.
type InstanceDetailContainer struct {
	Instance InstanceDetail `json:"instance"`
}

// VolumeSize has properties to set the volume of a database instance.
type VolumeSize struct {
	Size int `json:"size"`
}

// CreateInstanceParameters has properties for creating a trove database instance.
type CreateInstanceParameters struct {
	Databases         []CreateDatabaseParameters `json:"databases,omitempty"`
	FlavorRef         string                     `json:"flavorRef"`
	Name              string                     `json:"name,omitempty"`
	Users             []UserParameter            `json:"users,omitempty"`
	Volume            *VolumeSize                `json:"volume,omitempty"`
	NetworkInterfaces []NetworkInterface         `json:"nics,omitempty"`
}

// NetworkInterface represents a network interface to Trove
type NetworkInterface struct {
	NetID string `json:"net-id"`
}

// RestoreBackupParameters has properties for restoring a backup to a trove database instance.
type RestoreBackupParameters struct {
	FlavorRef    string       `json:"flavorRef"`
	RestorePoint RestorePoint `json:"restorePoint"`
	Name         string       `json:"name,omitempty"`
	Volume       *VolumeSize  `json:"volume,omitempty"`
}

// Instances will issue a GET request to retrieve all database instances.
func (databaseService Service) Instances() ([]Instance, error) {
	var container = InstancesContainer{}
	reqURL, err := databaseService.buildRequestURL("/instances")
	if err != nil {
		return container.Instances, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)

	return container.Instances, err
}

// Instance will issue a GET request to retrieve the specified database instance.
func (databaseService Service) Instance(instanceID string) (InstanceDetail, error) {
	var container = InstanceDetailContainer{}
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID)
	if err != nil {
		return container.Instance, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)
	return container.Instance, err
}

// CreateInstance will issue a POST request to create the specified database instance.
func (databaseService Service) CreateInstance(params CreateInstanceParameters) (InstanceDetail, error) {
	var inContainer = createInstanceParametersContainer{params}
	reqURL, err := databaseService.buildRequestURL("/instances")
	if err != nil {
		return InstanceDetail{}, err
	}

	var outContainer = InstanceDetailContainer{}
	err = misc.PostJSON(reqURL, databaseService.authenticator, &inContainer, &outContainer)
	return outContainer.Instance, err
}

// RestoreBackupInstance will issue a POST request to create a new Backup .
func (databaseService Service) RestoreBackupInstance(params RestoreBackupParameters) (InstanceDetail, error) {
	var inContainer = restoreBackupInstanceContainer{params}
	reqURL, err := databaseService.buildRequestURL("/instances")
	if err != nil {
		return InstanceDetail{}, err
	}

	var outContainer = InstanceDetailContainer{}
	err = misc.PostJSON(reqURL, databaseService.authenticator, &inContainer, &outContainer)
	return outContainer.Instance, err
}

// DeleteInstance will issue a delete query to delete the specified database instance.
func (databaseService Service) DeleteInstance(instanceID string) (err error) {
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, databaseService.authenticator)
}

// EnableRoot will issue a delete query to delete the specified database instance.
func (databaseService Service) EnableRoot(instanceID string) (UserParameter, error) {
	reqURL, err := databaseService.buildRequestURL("/instances/", instanceID, "/root")
	if err != nil {
		return UserParameter{}, err
	}

	var outContainer = userParameterContainer{}
	err = misc.PostJSON(reqURL, databaseService.authenticator, nil, &outContainer)
	return outContainer.User, err
}

type createInstanceParametersContainer struct {
	Instance CreateInstanceParameters `json:"instance"`
}

type restoreBackupInstanceContainer struct {
	Backup RestoreBackupParameters `json:"instance"`
}
