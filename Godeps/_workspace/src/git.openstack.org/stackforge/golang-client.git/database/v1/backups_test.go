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

import (
	"encoding/json"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleBackup = Backup{
	Created:     createdTime,
	Description: "Description",
	ID:          "agag",
	InstanceID:  "2t892t",
	LocationRef: "53225",
	Name:        "Name",
	ParentID:    "parentID",
	Size:        30,
	Status:      "Active",
	Updated:     updatedTime,
}

var sampleBackupBytes, _ = json.Marshal(sampleBackup)
var sampleBackupResponse = `{ "backup":` + string(sampleBackupBytes) + `}`
var sampleBackupsResponse = `{ "backups":[` + string(sampleBackupBytes) + `]}`

func TestBackup(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleBackupResponse, "/backups/56256")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	backup, err := service.Backup("56256")

	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleBackup, backup)
}

func TestBackups(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleBackupsResponse, "/backups")
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	backups, err := service.Backups()

	testUtil.IsNil(t, err)
	testUtil.Equals(t, 1, len(backups))
	testUtil.Equals(t, sampleBackup, backups[0])
}

func TestDeleteBackup(t *testing.T) {
	name := "user"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/backups/"+name)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	err := service.DeleteBackup(name)
	testUtil.IsNil(t, err)
}

func TestCreateBackup(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleBackupResponse, "/backups",
		`{"backup":{"name":"MyInstance1Backup","instance":"InstanceID","description":"My first backup"}}`)
	defer apiServer.Close()

	service := CreateDatabaseService(apiServer.URL)
	backup, err := service.CreateBackup(CreateBackupParameters{Name: "MyInstance1Backup", InstanceID: "InstanceID", Description: "My first backup"})

	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleBackup, backup)
}
