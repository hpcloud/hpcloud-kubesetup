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

// image.go
package image_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	image "git.openstack.org/stackforge/golang-client.git/image/v1"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var tokn = "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"

func TestListImages(t *testing.T) {
	anon := func(imageService *image.Service) {
		images, err := imageService.Images()
		if err != nil {
			t.Error(err)
		}

		if len(images) != 3 {
			t.Error(errors.New("Incorrect number of images found"))
		}
		expectedImage := image.Response{
			Name:            "Ubuntu Server 14.04.1 LTS (amd64 20140927) - Partner Image",
			ContainerFormat: "bare",
			DiskFormat:      "qcow2",
			CheckSum:        "6798a7d67ff0b241b6fe165798914d86",
			ID:              "bec3cab5-4722-40b9-a78a-3489218e22fe",
			Size:            255525376}
		// Verify first one matches expected values
		testUtil.Equals(t, expectedImage, images[0])
	}

	testImageServiceAction(t, "images?limit=21", sampleImagesData, anon)
}

func TestListImageDetails(t *testing.T) {
	anon := func(imageService *image.Service) {
		images, err := imageService.ImagesDetail()
		if err != nil {
			t.Error(err)
		}

		if len(images) != 2 {
			t.Error(errors.New("Incorrect number of images found"))
		}
		createdAt, _ := misc.NewDateTimeFromString(`"2014-09-29T14:44:31"`)
		updatedAt, _ := misc.NewDateTimeFromString(`"2014-09-29T15:33:37"`)
		owner := "10014302369510"
		virtualSize := int64(2525125)
		expectedImageDetail := image.DetailResponse{
			Status:          "active",
			Name:            "Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
			Deleted:         false,
			ContainerFormat: "bare",
			CreatedAt:       createdAt,
			DiskFormat:      "qcow2",
			UpdatedAt:       updatedAt,
			MinDisk:         8,
			Protected:       false,
			ID:              "8ca068c5-6fde-4701-bab8-322b3e7c8d81",
			MinRAM:          0,
			CheckSum:        "de1831ea85702599a27e7e63a9a444c3",
			Owner:           &owner,
			IsPublic:        true,
			DeletedAt:       nil,
			Properties: map[string]string{
				"com.ubuntu.cloud__1__milestone":    "release",
				"com.hp__1__os_distro":              "com.ubuntu",
				"description":                       "Ubuntu Server 12.04.5 LTS (amd64 20140927) for HP Public Cloud. Ubuntu Server is the world's most popular Linux for cloud environments. Updates and patches for Ubuntu 12.04.5 LTS will be available until 2017-04-26. Ubuntu Server is the perfect platform for all workloads from web applications to NoSQL databases and Hadoop. More information regarding Ubuntu Cloud is available from http://www.ubuntu.com/cloud and instructions for using Juju to deploy workloads are available from http://juju.ubuntu.com EULA: http://www.ubuntu.com/about/about-ubuntu/licensing Privacy Policy: http://www.ubuntu.com/privacy-policy",
				"com.ubuntu.cloud__1__suite":        "precise",
				"com.ubuntu.cloud__1__serial":       "20140927",
				"com.hp__1__bootable_volume":        "True",
				"com.hp__1__vendor":                 "Canonical",
				"com.hp__1__image_lifecycle":        "active",
				"com.hp__1__image_type":             "disk",
				"os_version":                        "12.04",
				"architecture":                      "x86_64",
				"os_type":                           "linux-ext4",
				"com.ubuntu.cloud__1__stream":       "server",
				"com.ubuntu.cloud__1__official":     "True",
				"com.ubuntu.cloud__1__published_at": "2014-09-29T15:33:36"},
			Size:        261423616,
			VirtualSize: &virtualSize}
		testUtil.Equals(t, expectedImageDetail, images[0])
	}

	testImageServiceAction(t, "images/detail?limit=21", sampleImageDetailsData, anon)
}

func TestNameFilterUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=21&name=CentOS+deprecated",
		image.QueryParameters{Name: "CentOS deprecated"})
}

func TestStatusUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=21&status=active",
		image.QueryParameters{Status: "active"})
}

func TestMinMaxSizeUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=21&size_max=5300014&size_min=100158",
		image.QueryParameters{MinSize: 100158, MaxSize: 5300014})
}

func TestMarkerLimitUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=20&marker=bec3cab5-4722-40b9-a78a-3489218e22fe",
		image.QueryParameters{Marker: "bec3cab5-4722-40b9-a78a-3489218e22fe", Limit: 20})
}

func TestContainerFormatFilterUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?container_format=bare&limit=21",
		image.QueryParameters{ContainerFormat: "bare"})
}

func TestSortKeySortUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=21&sort_key=id",
		image.QueryParameters{SortKey: "id"})
}

func TestSortDirSortUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=21&sort_dir=asc",
		image.QueryParameters{SortDirection: image.Asc})
}

func testImageQueryParameter(t *testing.T, uriEndsWith string, queryParameters image.QueryParameters) {
	anon := func(imageService *image.Service) {
		_, _ = imageService.QueryImages(queryParameters)
	}

	testImageServiceAction(t, uriEndsWith, sampleImagesData, anon)
}

func testImageServiceAction(t *testing.T, uriEndsWith string, testData string, imageServiceAction func(*image.Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect URL created, expected:" + uriEndsWith + " at the end, actual URL:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	imageService := image.NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: apiServer.URL})
	imageServiceAction(&imageService)
}

func testImageServiceGetAction(t *testing.T, uriEndsWith string, testHeaders http.Header, testData []byte, imageServiceAction func(*image.Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect URL created, expected:" + uriEndsWith + " at the end, actual URL:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetTestRequestServer(t, tokn, 200, testHeaders, testData, anon)
	defer apiServer.Close()

	imageService := image.NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: apiServer.URL})
	imageServiceAction(&imageService)
}

func testImageServiceHeadAction(t *testing.T, uriEndsWith string, testHeaders http.Header, testData []byte, imageServiceAction func(*image.Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateHeadTestRequestServer(t, tokn, 200, testHeaders, testData, anon)
	defer apiServer.Close()

	imageService := image.NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: apiServer.URL})
	imageServiceAction(&imageService)
}

func TestImageDetail(t *testing.T) {
	testHeaders := http.Header{
		"X-Image-Meta-property-purpose": []string{"test"},
		"X-Image-Meta-status":           []string{"active"},
		"X-Image-Meta-owner":            []string{"54026737306152"},
		"X-Image-Meta-name":             []string{"Ubuntu Minimal"},
		"X-Image-Meta-container_format": []string{"bare"},
		"X-Image-Meta-created_at":       []string{"2015-04-03T18:19:40"},
		"X-Image-Meta-min_ram":          []string{"4096"},
		"X-Image-Meta-updated_at":       []string{"2015-04-03T18:19:42"},
		"X-Image-Meta-property-os":      []string{"ubuntu"},
		"X-Image-Meta-id":               []string{"84cf82c0-41e7-4f5f-8722-40b03878b4a9"},
		"X-Image-Meta-deleted":          []string{"False"},
		"X-Image-Meta-checksum":         []string{"4f783f3917ed4c663c9a983f4ee046fc"},
		"X-Image-Meta-protected":        []string{"False"},
		"X-Image-Meta-min_disk":         []string{"10"},
		"X-Image-Meta-size":             []string{"40894464"},
		"X-Image-Meta-is_public":        []string{"False"},
		"X-Image-Meta-disk_format":      []string{"iso"},
	}

	testBody := []byte("")

	anon := func(imageService *image.Service) {
		detail, err := imageService.ImageDetail("84cf82c0-41e7-4f5f-8722-40b03878b4a9")
		if err != nil {
			t.Error(err)
		}

		createdAt, err := misc.NewDateTimeFromString("2015-04-03T18:19:40")
		if err != nil {
			t.Error(err)
		}

		updatedAt, err := misc.NewDateTimeFromString("2015-04-03T18:19:42")
		if err != nil {
			t.Error(err)
		}

		owner := "54026737306152"

		expectedDetails := image.DetailResponse{
			CheckSum:        "4f783f3917ed4c663c9a983f4ee046fc",
			ContainerFormat: "bare",
			CreatedAt:       createdAt,
			Deleted:         false,
			DeletedAt:       nil,
			DiskFormat:      "iso",
			ID:              "84cf82c0-41e7-4f5f-8722-40b03878b4a9",
			IsPublic:        false,
			MinDisk:         10,
			MinRAM:          4096,
			Name:            "Ubuntu Minimal",
			Owner:           &owner,
			UpdatedAt:       updatedAt,
			Properties: map[string]string{
				"purpose": "test",
				"os":      "ubuntu",
			},
			Protected:   false,
			Status:      "active",
			Size:        40894464,
			VirtualSize: nil,
		}

		testUtil.Equals(t, expectedDetails, detail)
	}

	testImageServiceHeadAction(t, "images/84cf82c0-41e7-4f5f-8722-40b03878b4a9", testHeaders, testBody, anon)
}

var sampleImagesData = `{
   "images":[
      {
         "name":"Ubuntu Server 14.04.1 LTS (amd64 20140927) - Partner Image",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"6798a7d67ff0b241b6fe165798914d86",
         "id":"bec3cab5-4722-40b9-a78a-3489218e22fe",
         "size":255525376
      },
      {
         "name":"Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"de1831ea85702599a27e7e63a9a444c3",
         "id":"8ca068c5-6fde-4701-bab8-322b3e7c8d81",
         "size":261423616
      },
      {
         "name":"HP_LR-PC_Load_Generator_12-02_Windows-2008R2x64",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"052d70c2b4d4988a8816197381e9083a",
         "id":"12b9c19b-8823-4f40-9531-0f05fb0933f2",
         "size":14012055552
      }
   ]
}`

var sampleImageDetailsData = `{
   "images":[
      {
         "status":"active",
         "name":"Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
         "deleted":false,
         "container_format":"bare",
         "created_at":"2014-09-29T14:44:31",
         "disk_format":"qcow2",
         "updated_at":"2014-09-29T15:33:37",
         "min_disk":8,
         "protected":false,
         "id":"8ca068c5-6fde-4701-bab8-322b3e7c8d81",
         "min_ram":0,
         "checksum":"de1831ea85702599a27e7e63a9a444c3",
         "owner":"10014302369510",
         "is_public":true,
         "deleted_at":null,
         "properties":{
            "com.ubuntu.cloud__1__milestone":"release",
            "com.hp__1__os_distro":"com.ubuntu",
            "description":"Ubuntu Server 12.04.5 LTS (amd64 20140927) for HP Public Cloud. Ubuntu Server is the world's most popular Linux for cloud environments. Updates and patches for Ubuntu 12.04.5 LTS will be available until 2017-04-26. Ubuntu Server is the perfect platform for all workloads from web applications to NoSQL databases and Hadoop. More information regarding Ubuntu Cloud is available from http://www.ubuntu.com/cloud and instructions for using Juju to deploy workloads are available from http://juju.ubuntu.com EULA: http://www.ubuntu.com/about/about-ubuntu/licensing Privacy Policy: http://www.ubuntu.com/privacy-policy",
            "com.ubuntu.cloud__1__suite":"precise",
            "com.ubuntu.cloud__1__serial":"20140927",
            "com.hp__1__bootable_volume":"True",
            "com.hp__1__vendor":"Canonical",
            "com.hp__1__image_lifecycle":"active",
            "com.hp__1__image_type":"disk",
            "os_version":"12.04",
            "architecture":"x86_64",
            "os_type":"linux-ext4",
            "com.ubuntu.cloud__1__stream":"server",
            "com.ubuntu.cloud__1__official":"True",
            "com.ubuntu.cloud__1__published_at":"2014-09-29T15:33:36"
         },
         "size":261423616,
		 "virtual_size":2525125
      },
      {
         "status":"active",
         "name":"Windows Server 2008 Enterprise SP2 x64 Volume License 20140415 (b)",
         "deleted":false,
         "container_format":"bare",
         "created_at":"2014-04-25T19:53:24",
         "disk_format":"qcow2",
         "updated_at":"2014-04-25T19:57:11",
         "min_disk":30,
         "protected":true,
         "id":"1294610e-fdc4-579b-829b-d0c9f5c0a612",
         "min_ram":0,
         "checksum":"37208aa6d49929f12132235c5f834f2d",
         "owner":null,
         "is_public":true,
         "deleted_at":null,
         "properties":{
            "hp_image_license":"1002",
            "com.hp__1__os_distro":"com.microsoft.server",
            "com.hp__1__image_lifecycle":"active",
            "com.hp__1__image_type":"disk",
            "architecture":"x86_64",
            "com.hp__1__license_os":"1002",
            "com.hp__1__bootable_volume":"true"
         },
         "size":6932856832,
		"virtual_size":null
      }
   ]
}`

func TestPrivateCloudImages(t *testing.T) {

	// Read private cloud response content from file
	imagesTestFilePath := "./testdata/private_cloud_images.json"
	imagesTestFileContent, err := ioutil.ReadFile(imagesTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to read JSON file %s: '%s'", imagesTestFilePath, err.Error()))
	}

	// Decode the content
	sampleImgs := imagesResponse{}
	err = json.Unmarshal(imagesTestFileContent, &sampleImgs)
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", imagesTestFilePath, err.Error()))
	}

	// Test the SDK API imageService.Images()
	anon := func(imageService *image.Service) {
		images, err := imageService.Images()
		if err != nil {
			t.Error(err)
		}

		if len(images) != len(sampleImgs.Images) {
			t.Error(errors.New("Incorrect number of images found"))
		}

		// Verify returned images match original sample images
		testUtil.Equals(t, sampleImgs.Images, images)
	}

	testImageServiceAction(t, "images?limit=21", string(imagesTestFileContent), anon)
}

func TestPrivateCloudImagesDetail(t *testing.T) {
	// Read private cloud response content from file
	imagesTestFilePath := "./testdata/private_cloud_images_detail.json"
	imagesTestFileContent, err := ioutil.ReadFile(imagesTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to read JSON file %s: '%s'", imagesTestFilePath, err.Error()))
	}

	// Decode the content
	sampleImgs := imagesDetailResponse{}
	err = json.Unmarshal(imagesTestFileContent, &sampleImgs)
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", imagesTestFilePath, err.Error()))
	}

	// Test the SDK API imageService.ImagesDetail()
	anon := func(imageService *image.Service) {
		images, err := imageService.ImagesDetail()
		if err != nil {
			t.Error(err)
		}

		if len(images) != len(sampleImgs.Images) {
			t.Error(errors.New("Incorrect number of images found"))
		}

		// Verify returned images match original sample images
		testUtil.Equals(t, sampleImgs.Images, images)
	}

	testImageServiceAction(t, "images/detail?limit=21", string(imagesTestFileContent), anon)
}

type imagesResponse struct {
	Images []image.Response `json:"images"`
}

type imagesDetailResponse struct {
	Images []image.DetailResponse `json:"images"`
}
