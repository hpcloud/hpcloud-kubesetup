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

package objectstorage_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	objectstorage "git.openstack.org/stackforge/golang-client.git/objectstorage/v1"
)

// TestObjectStorageAPI needs to be broken up into individual tests.
func TestObjectStorageAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	container := os.Getenv("TEST_OBJECTSTORAGE_CONTAINERNAME")
	if container == "" {
		t.Skip("No container specified for integration test in TEST_OBJECTSTORAGE_CONTAINERNAME env variable.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	tokenID, err := authenticator.GetToken()
	if err != nil {
		t.Fatal("There was an error authenticating:", err)
	}

	url, err := authenticator.GetServiceURL("object-store", "1.0")
	if err != nil {
		t.Fatal("There was an error determining the object-store service url:", err)
	}

	hdr, err := objectstorage.GetAccountMeta(url, tokenID)
	if err != nil {
		t.Fatal("There was an error getting account metadata:", err)
	}

	// Create a new container.
	if err = objectstorage.PutContainer(url+container, tokenID,
		"X-Log-Retention", "true"); err != nil {
		t.Fatal("PutContainer Error:", err)
	}

	// Get a list of all the containers at the selected endoint.
	containersJSON, err := objectstorage.ListContainers(0, "", url, tokenID)
	if err != nil {
		t.Fatal(err)
	}

	type containerType struct {
		Name         string
		Bytes, Count int
	}
	containersList := []containerType{}

	if err = json.Unmarshal(containersJSON, &containersList); err != nil {
		t.Fatal(err)
	}

	found := false
	for i := 0; i < len(containersList); i++ {
		if containersList[i].Name == container {
			found = true
		}
	}
	if !found {
		t.Fatal("Created container is missing from downloaded containersList")
	}

	// Set and Get container metadata.
	if err = objectstorage.SetContainerMeta(url+container, tokenID,
		"X-Container-Meta-fubar", "false"); err != nil {
		t.Fatal(err)
	}

	hdr, err = objectstorage.GetContainerMeta(url+container, tokenID)
	if err != nil {
		t.Fatal(fmt.Sprint("GetContainerMeta Error:", err))
	}
	if hdr.Get("X-Container-Meta-fubar") != "false" {
		t.Fatal("container meta does not match")
	}

	// Create an object in a container.
	var fContent []byte
	srcFile := "10-objectstore.go"
	fContent, err = ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	object := container + "/" + srcFile
	if err = objectstorage.PutObject(&fContent, url+object, tokenID,
		"X-Object-Meta-fubar", "false"); err != nil {
		t.Fatal(err)
	}
	objectsJSON, err := objectstorage.ListObjects(0, "", "", "", "",
		url+container, tokenID)

	type objectType struct {
		Name, Hash, Content_type, Last_modified string
		Bytes                                   int
	}
	objectsList := []objectType{}

	if err = json.Unmarshal(objectsJSON, &objectsList); err != nil {
		t.Fatal(err)
	}
	found = false
	for i := 0; i < len(objectsList); i++ {
		if objectsList[i].Name == srcFile {
			found = true
		}
	}
	if !found {
		t.Fatal("created object is missing from the objectsList")
	}

	// Manage object metadata
	if err = objectstorage.SetObjectMeta(url+object, tokenID,
		"X-Object-Meta-fubar", "true"); err != nil {
		t.Fatal("SetObjectMeta Error:", err)
	}
	hdr, err = objectstorage.GetObjectMeta(url+object, tokenID)
	if err != nil {
		t.Fatal("GetObjectMeta Error:", err)
	}
	if hdr.Get("X-Object-Meta-fubar") != "true" {
		t.Fatal("SetObjectMeta Error:", err)
	}

	// Retrieve an object and check that it is the same as what as uploaded.
	_, body, err := objectstorage.GetObject(url+object, tokenID)
	if err != nil {
		t.Fatal("GetObject Error:", err)
	}
	if !bytes.Equal(fContent, body) {
		t.Fatal("GetObject Error:", "byte comparison of uploaded != downloaded")
	}

	// Duplication (Copy) an existing object.
	if err = objectstorage.CopyObject(url+object, "/"+object+".dup", tokenID); err != nil {
		t.Fatal("CopyObject Error:", err)
	}

	// Delete the objects.
	if err = objectstorage.DeleteObject(url+object, tokenID); err != nil {
		t.Fatal("DeleteObject Error:", err)
	}
	if err = objectstorage.DeleteObject(url+object+".dup", tokenID); err != nil {
		t.Fatal("DeleteObject Error:", err)
	}

	// Delete the container that was previously created.
	if err = objectstorage.DeleteContainer(url+container, tokenID); err != nil {
		t.Fatal("DeleteContainer Error:", err)
	}
}
