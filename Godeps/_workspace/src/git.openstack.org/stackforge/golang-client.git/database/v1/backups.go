// Copyright (c) 2015 Hewlett-Packard Development Company, L.P.
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

import "git.openstack.org/stackforge/golang-client.git/misc"

// RestorePoint has the reference to the backup to restore to.
type RestorePoint struct {
	BackupRef string `json:"backupRef"`
}

// CreateBackupParameters has properties for creating a backup of a trove database instance.
type CreateBackupParameters struct {
	Name        string `json:"name,omitempty"`
	InstanceID  string `json:"instance"`
	Description string `json:"description,omitempty"`
}

// Backup contains properties of an backed up instance.
type Backup struct {
	Created     misc.RFC8601DateTime `json:"created,omitempty"`
	Description string               `json:"description,omitempty"`
	Datastore   Datastore            `json:"datastore,omitempty"`
	ID          string               `json:"id,omitempty"`
	InstanceID  string               `json:"instance_id"`
	LocationRef string               `json:"locationRef,omitempty"`
	Name        string               `json:"name"`
	ParentID    string               `json:"parent_id,omitempty"`
	Size        int                  `json:"size,omitempty"`
	Status      string               `json:"status,omitempty"`
	Updated     misc.RFC8601DateTime `json:"updated,omitempty"`
}

// Backups will issue a GET request to retrieve all backups.
func (databaseService Service) Backups() ([]Backup, error) {
	var container = backupsContainer{}
	reqURL, err := databaseService.buildRequestURL("/backups")
	if err != nil {
		return container.Backups, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)

	return container.Backups, err
}

// Backup will issue a GET request to retrieve the specified backup.
func (databaseService Service) Backup(backupID string) (Backup, error) {

	var container = backupContainer{}
	reqURL, err := databaseService.buildRequestURL("/backups/", backupID)
	if err != nil {
		return container.Backup, err
	}

	err = misc.GetJSON(reqURL, databaseService.authenticator, &container)
	return container.Backup, err
}

// CreateBackup will issue a POST request to create a new Backup.
func (databaseService Service) CreateBackup(params CreateBackupParameters) (Backup, error) {

	var inContainer = createBackupContainer{params}
	reqURL, err := databaseService.buildRequestURL("/backups")
	if err != nil {
		return Backup{}, err
	}

	var outContainer = backupContainer{}
	err = misc.PostJSON(reqURL, databaseService.authenticator, &inContainer, &outContainer)
	return outContainer.Backup, err
}

// DeleteBackup will issue a delete query to delete the specified backup.
func (databaseService Service) DeleteBackup(backupID string) (err error) {
	reqURL, err := databaseService.buildRequestURL("/backups/", backupID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, databaseService.authenticator)
}

type createBackupContainer struct {
	Backup CreateBackupParameters `json:"backup"`
}

type backupContainer struct {
	Backup Backup `json:"backup"`
}

type backupsContainer struct {
	Backups []Backup `json:"backups"`
}
