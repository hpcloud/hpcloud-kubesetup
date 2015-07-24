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

package compute

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleTime = time.Now().In(time.UTC)

var createServerResponse = CreateServerResponse{
	ID:             "3b451a21-c044-4981-8bea-70100bc6cf4a",
	Links:          []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}},
	AdminPass:      "ag4awrg",
	SecurityGroups: []SecurityGroup{SecurityGroup{Name: "default"}},
}

var serverDetail = ServerDetail{
	ID:               "my_server_id_1",
	Name:             "my_server_name_1",
	Status:           "my_server_status_1",
	Created:          sampleTime,
	Updated:          &sampleTime,
	HostID:           "my_server_hostId_1",
	Addresses:        make(map[string][]Address),
	Links:            []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}},
	Image:            ImageWrapper{Image: &Image{ID: "image_id1", Links: []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}}}},
	Flavor:           Flavor{ID: "image_id1", Links: []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}}},
	TaskState:        "my_OS-EXT-STS_task_state_1",
	VMState:          "my_OS-EXT-STS_vm_state_1",
	PowerState:       1,
	AvailabilityZone: "my_zone_a",
	UserID:           "my_user_id_1",
	TenantID:         "my_tenant_id_1",
	AccessIPv4:       "192.168.0.12",
	AccessIPv6:       "my_accessIPv6_1",
	ConfigDrive:      "my_config_drive_1",
	Progress:         2,
	MetaData:         make(map[string]string),
	AdminPass:        "my_adminPass_1",
	KeyName:          "my_key_name",
}

var servDetailContainer = serversDetailContainer{ServersDetail: []ServerDetail{serverDetail}}

var serverDetailPayloadBytes, _ = json.Marshal(serverDetail)
var serversDetailJSONPayload = `{ "servers": [` + string(serverDetailPayloadBytes) + `]}`
var serverDetailJSONPayload = `{ "server": ` + string(serverDetailPayloadBytes) + `}`

var createServerPayloadBytes, _ = json.Marshal(createServerResponse)
var createServerJSONPayload = `{ "server": ` + string(createServerPayloadBytes) + `}`

var tokn = "test-token-1"

func TestQueryServersDetail(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, serversDetailJSONPayload, "/servers/detail?limit=21&name=foo")
	defer apiServer.Close()

	computeService := CreateComputeService(apiServer.URL)
	queryParameters := ServerDetailQueryParameters{Name: "foo"}
	servers, err := computeService.QueryServersDetail(queryParameters)

	testUtil.IsNil(t, err)
	testUtil.Assert(t, len(servers) == 1, "Expected 1 server.")
	testUtil.Equals(t, serverDetail, servers[0])
}

func TestCreateServer(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, createServerJSONPayload, "servers",
		`{"server":{"name":"my_server","imageRef":"8c3cd338-1282-4fbb-bbaf-2256ff97c7b7","key_name":"my_key_name","flavorRef":"101","maxcount":1,"mincount":1,"user_data":"my_user_data","availability_zone":"az1","networks":[{"uuid":"1111d337-0282-4fbb-bbaf-2256ff97c7b7","port":"881"}],"security_groups":[{"name":"my_security_group_123"}]}}`)
	defer apiServer.Close()

	serverService := CreateComputeService(apiServer.URL)
	az := "az1"
	userData := "my_user_data"
	maxMinCount := int32(1)
	var serverCreationParameters = ServerCreationParameters{
		Name:             "my_server",
		ImageRef:         "8c3cd338-1282-4fbb-bbaf-2256ff97c7b7",
		KeyPairName:      "my_key_name",
		FlavorRef:        "101",
		MaxCount:         &maxMinCount,
		MinCount:         &maxMinCount,
		AvailabilityZone: &az,
		UserData:         &userData,
		Networks:         []ServerNetworkParameters{ServerNetworkParameters{UUID: "1111d337-0282-4fbb-bbaf-2256ff97c7b7", Port: "881"}},
		SecurityGroups:   []SecurityGroup{SecurityGroup{Name: "my_security_group_123"}}}

	result, err := serverService.CreateServer(serverCreationParameters)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, createServerResponse, result)
}

func TestDeleteServer(t *testing.T) {

	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/servers/server")
	defer apiServer.Close()

	serverService := CreateComputeService(apiServer.URL)

	err := serverService.DeleteServer("server")
	testUtil.IsNil(t, err)
}

func TestServerAction(t *testing.T) {

	marshaledserverInfoDetail, err := json.Marshal(serverDetail)
	if err != nil {
		t.Error(err)
	}

	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				w.Header().Set("Version", "my_server_version.1.0.0")
				w.WriteHeader(201)
				w.Write([]byte(marshaledserverInfoDetail))
				return
			}
			t.Error(errors.New("Failed: r.Method == POST"))
		}))
	defer apiServer.Close()

	serverService := CreateComputeService(apiServer.URL)

	err = serverService.ServerAction("server1", "action1", "key1", "value1")
	testUtil.IsNil(t, err)
}

func TestServers(t *testing.T) {

	var servers = []Server{
		{
			ID:    "123cd338-1282-4fbb-bbaf-2256ff97c111",
			Name:  "my_server_1",
			Links: []Link{{"href_11", "rel_11"}, {"href_12", "rel_12"}}},
		{
			ID:    "123cd338-1282-4fbb-bbaf-2256ff97c222",
			Name:  "my_server_2",
			Links: []Link{{"href_21", "rel_21"}, {"href_22", "rel_22"}}},
		{
			ID:    "123cd338-1282-4fbb-bbaf-2256ff97c333",
			Name:  "my_server_3",
			Links: []Link{{"href_31", "rel_31"}, {"href_32", "rel_32"}}}}

	var container = serversContainer{Servers: servers}

	marshaledContainer, err := json.Marshal(container)
	if err != nil {
		t.Error(err)
	}

	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Version", "my_server_version.1.0.0")
				w.WriteHeader(200)
				w.Write([]byte(marshaledContainer))
				return
			}
			t.Error(errors.New("Failed: r.Method == GET"))
		}))
	defer apiServer.Close()

	serverService := CreateComputeService(apiServer.URL)

	result, err := serverService.Servers()
	testUtil.IsNil(t, err)

	testUtil.Equals(t, servers, result)
}

func TestServerDetails(t *testing.T) {

	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Version", "my_server_version.1.0.0")
				w.WriteHeader(200)
				w.Write([]byte(serversDetailJSONPayload))
				return
			}
			t.Error(errors.New("Failed: r.Method == GET"))
		}))
	defer apiServer.Close()

	serverService := CreateComputeService(apiServer.URL)

	result, err := serverService.ServerDetails()
	testUtil.IsNil(t, err)

	testUtil.Equals(t, servDetailContainer.ServersDetail, result)
}

func TestServerDetail(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, serverDetailJSONPayload, "/servers/my_server")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	result, err := service.ServerDetail("my_server")
	testUtil.IsNil(t, err)

	testUtil.Equals(t, serverDetail, result)
}

var testServeID = "125311"

func TestGetServerMetadata(t *testing.T) {
	sampleMetadata := map[string]string{
		"status":  "OK",
		"message": "Processing is working!"}

	mockResponseObject := serverMetadataContainer{Metadata: sampleMetadata}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "servers/"+testServeID+"/metadata")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	foundMetadata, err := service.ServerMetadata(testServeID)
	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleMetadata, foundMetadata)
}

func TestGetServerMetaItem(t *testing.T) {
	sampleMetadata := map[string]string{"status": "OK"}

	mockResponseObject := serverMetaItemContainer{Meta: sampleMetadata}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "servers/"+testServeID+"/metadata/status")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	foundItem, err := service.ServerMetadataItem(testServeID, "status")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "OK", foundItem)
}

func TestSetServerMetaItem(t *testing.T) {
	sampleMetadata := map[string]string{"item": "value"}
	apiServer := CreateTestRequestServer(t, "PUT", tokn, "/servers/"+testServeID+"/metadata/item", `{"meta":{"item":"value"}}`, serverMetaItemContainer{Meta: sampleMetadata})
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	actualMetadata, err := service.SetServerMetadataItem(testServeID, "item", "value")
	if err != nil {
		t.Error(err)
	}
	testUtil.Equals(t, sampleMetadata, actualMetadata)
}

func TestSetServerMetadata(t *testing.T) {
	sampleMetadata := map[string]string{"item": "value"}
	apiServer := CreateTestRequestServer(t, "POST", tokn, "/servers/"+testServeID+"/metadata", `{"metadata":{"item":"value"}}`, serverMetadataContainer{Metadata: sampleMetadata})
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	actualMetadata, err := service.SetServerMetadata(testServeID, sampleMetadata)
	if err != nil {
		t.Error(err)
	}
	testUtil.Equals(t, sampleMetadata, actualMetadata)
}

func TestDeleteServerMetadata(t *testing.T) {
	apiServer := CreateTestRequestServer(t, "POST", tokn, "/servers/"+testServeID+"/metadata", `{"metadata":{}}`, make(map[string]string, 0))
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	err := service.DeleteServerMetadata(testServeID)
	if err != nil {
		t.Error(err)
	}
}

func CreateComputeService(url string) Service {
	return NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: url})
}

func CreateTestRequestServer(t *testing.T, expectedMethod string, expectedAuthTokenValue string, urlEndsWith string, expectedRequestBody string, responseOutput interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testUtil.HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			requestBytesArr, _ := ioutil.ReadAll(r.Body)
			testUtil.Equals(t, expectedRequestBody, string(requestBytesArr))
			reqURL := r.URL.String()
			if !strings.HasSuffix(reqURL, urlEndsWith) {
				t.Error(errors.New("Incorrect URL created, expected '" + urlEndsWith + "' at the end, actual URL:" + reqURL))
			}
			if r.Method == expectedMethod {
				data, _ := json.Marshal(responseOutput)
				w.Write(data)
				w.WriteHeader(200)
				return
			}

			t.Error(errors.New("Failed: r.Method == " + expectedMethod))
		}))
}

// Tests for the special case where image is "" instead of an image structure

var staticCreatedRFC8601, _ = misc.NewDateTimeFromString("2015-04-03T22:47:44")
var staticCreatedTime = staticCreatedRFC8601.Time()

var serverDetailNoImage = ServerDetail{
	ID:               "my_server_id_1",
	Name:             "my_server_name_1",
	Status:           "my_server_status_1",
	Created:          staticCreatedTime,
	Updated:          &staticCreatedTime,
	HostID:           "my_server_hostId_1",
	Addresses:        make(map[string][]Address),
	Links:            []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}},
	Image:            ImageWrapper{},
	Flavor:           Flavor{ID: "image_id1", Links: []Link{{"href_1", "rel_1"}, {"href_2", "rel_2"}}},
	TaskState:        "my_OS-EXT-STS_task_state_1",
	VMState:          "my_OS-EXT-STS_vm_state_1",
	PowerState:       1,
	AvailabilityZone: "my_zone_a",
	UserID:           "my_user_id_1",
	TenantID:         "my_tenant_id_1",
	AccessIPv4:       "192.168.0.12",
	AccessIPv6:       "my_accessIPv6_1",
	ConfigDrive:      "my_config_drive_1",
	Progress:         2,
	MetaData:         make(map[string]string),
	AdminPass:        "my_adminPass_1",
	KeyName:          "my_key_name",
}

var serverDetailNoImagePayload = `{ "server":
    {
       "id":"my_server_id_1",
       "name":"my_server_name_1",
       "status":"my_server_status_1",
       "created":"2015-04-03T22:47:44Z",
       "updated":"2015-04-03T22:47:44Z",
       "hostId":"my_server_hostId_1",
       "addresses":{

       },
       "links":[
          {
             "href":"href_1",
             "rel":"rel_1"
          },
          {
             "href":"href_2",
             "rel":"rel_2"
          }
       ],
       "image":"",
       "flavor":{
          "id":"image_id1",
          "links":[
             {
                "href":"href_1",
                "rel":"rel_1"
             },
             {
                "href":"href_2",
                "rel":"rel_2"
             }
          ]
       },
       "OS-EXT-STS:task_state":"my_OS-EXT-STS_task_state_1",
       "OS-EXT-STS:vm_state":"my_OS-EXT-STS_vm_state_1",
       "OS-EXT-STS:power_state":1,
       "OS-EXT-AZ:availability_zone":"my_zone_a",
       "user_id":"my_user_id_1",
       "tenant_id":"my_tenant_id_1",
       "accessIPv4":"192.168.0.12",
       "accessIPv6":"my_accessIPv6_1",
       "config_drive":"my_config_drive_1",
       "progress":2,
       "metadata":{

       },
       "adminPass":"my_adminPass_1",
       "key_name":"my_key_name"
    }
}`

func TestServerDetailNoImage(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, serverDetailNoImagePayload, "/servers/my_server")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)

	result, err := service.ServerDetail("my_server")
	testUtil.IsNil(t, err)

	testUtil.Equals(t, serverDetailNoImage, result)
}
