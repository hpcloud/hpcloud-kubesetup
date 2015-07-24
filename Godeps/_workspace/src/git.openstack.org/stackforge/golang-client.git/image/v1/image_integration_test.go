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

package image_test

import (
	"testing"

	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	image "git.openstack.org/stackforge/golang-client.git/image/v1"
)

// Image examples.
func TestImageServiceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	imageService := image.NewService(authenticator)
	imagesDetails, err := imageService.ImagesDetail()
	if err != nil {
		t.Fatal("Cannot access images:", err)
	}

	var imageIDs = make([]string, 0)
	for _, element := range imagesDetails {
		imageIDs = append(imageIDs, element.ID)
	}

	if len(imageIDs) == 0 {
		t.Fatal("No images found, check to make sure access is correct")
	}
}
