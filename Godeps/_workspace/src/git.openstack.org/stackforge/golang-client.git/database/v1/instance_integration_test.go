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

package database_test

import (
	"os"
	"testing"
	"time"

	database "git.openstack.org/stackforge/golang-client.git/database/v1"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

// Image examples.
func TestDatabaseServiceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	if os.Getenv("OS_AUTH_URL") == "" {
		t.Skip("Cannot run integration test as the Openstack env vars aren't set.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	authenticator.(requester.Manager).SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	databaseService := database.NewService(authenticator)

	user := database.UserParameter{Name: "User1", Password: "ag92340gv"}
	db := database.CreateDatabaseParameters{Name: "newDb"}
	createInstanceParams := database.CreateInstanceParameters{
		Databases: []database.CreateDatabaseParameters{db},
		FlavorRef: "1001",
		Name:      "Instance1",
		Users:     []database.UserParameter{user},
	}

	instance, err := databaseService.CreateInstance(createInstanceParams)
	if err != nil {
		t.Fatal("Cannot create instance:", err)
	}

	WaitUntilActiveOrError(t, databaseService, instance.ID)

	foundInstance := false
	instances, err := databaseService.Instances()
	if err != nil {
		t.Fatal("Cannot query Instances:", err)
	}

	for _, i := range instances {
		if i.ID == instance.ID {
			foundInstance = true
			break
		}
	}

	if !foundInstance {
		t.Fatal("Cannot find new instance")
	}

	err = databaseService.CreateUser(instance.ID, database.UserParameter{Name: "username", Password: "39hfnw282"})
	if err != nil {
		t.Fatal("Cannot create user", err)
	}

	err = databaseService.DeleteUser(instance.ID, "username")
	if err != nil {
		t.Fatal("Cannot delete user", err)
	}

	pwd, err := databaseService.EnableRoot(instance.ID)
	if err != nil {
		t.Fatal("Cannot enable root", err)
	}
	testUtil.Assert(t, pwd.Name != "", "No Name")

	backup, err := databaseService.CreateBackup(database.CreateBackupParameters{Name: "NewBackup", InstanceID: instance.ID, Description: "Test Description"})
	if err != nil {
		t.Fatal("Cannot make backup", err)
	}
	WaitUntilActiveOrError(t, databaseService, instance.ID)

	err = databaseService.DeleteInstance(instance.ID)
	if err != nil {
		t.Fatal("Delete instance didn't work:", err)
	}

	newInstance, err := databaseService.RestoreBackupInstance(database.RestoreBackupParameters{Name: "Instance1", RestorePoint: database.RestorePoint{BackupRef: backup.ID}, FlavorRef: "1001"})
	if err != nil {
		t.Fatal("Cannot restore a backup", err)
	}

	WaitUntilActiveOrError(t, databaseService, newInstance.ID)

	err = databaseService.DeleteBackup(backup.ID)
	if err != nil {
		t.Fatal("Cannot delete backup", err)
	}

	err = databaseService.DeleteInstance(newInstance.ID)
	if err != nil {
		t.Fatal("Delete instance didn't work:", err)
	}
}

func WaitUntilActiveOrError(t *testing.T, databaseService database.Service, instanceID string) {
	for i := 10; i > 0; i-- {
		foundInstanceLookup, err := databaseService.Instance(instanceID)
		if err != nil {
			t.Fatal("Error looking up instance", err)
		}

		if foundInstanceLookup.Status == database.ActiveStatus {
			return
		}

		if foundInstanceLookup.Status == database.FailedStatus {
			t.Fatal("Instance in a Failed state.")
		}

		time.Sleep(time.Duration(30 * time.Second))
	}

	t.Fatal("Unable to wait for Active status.")
}
