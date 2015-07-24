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

// auth_test.go
package identity_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	common "git.openstack.org/stackforge/golang-client.git/identity/common"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	testutil "git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestExpiredTokenErrorIsReturned(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testutil.Equals(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(getSampleAuthPayload("2001-01-07T07:22:52.184Z")))
			return
		}))
	defer apiServer.Close()

	params := common.AuthenticationParameters{AuthURL: apiServer.URL, Username: "chris", Password: "Password", Region: "region-a.geo-1"}
	authenticator := identity.Authenticate(params)
	_, err := authenticator.GetToken()
	if err == nil {
		t.Fatal("Expected an error")
	}
	testutil.Equals(t, "Error: The auth token has an invalid expiration.", err.Error())
}

func TestTenantIDPrecedenceOverTenantNameIfBothIncluded(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"tenantId\":\"tenantID\",\"passwordCredentials\":{\"password\":\"Password\",\"username\":\"chris\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, Username: "chris", Password: "Password", Region: "region-a.geo-1", TenantID: "tenantID", TenantName: "TenantName"}
			return identity.Authenticate(params)
		})
}

func TestAuthenticateWithOnlyUsernameAndPasswordValid(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"passwordCredentials\":{\"password\":\"Password\",\"username\":\"chris\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, Username: "chris", Password: "Password", Region: "region-a.geo-1"}
			return identity.Authenticate(params)
		})
}

func TestAuthenticateWithOnlyUsernameAndPasswordAndTenantIDValid(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"tenantId\":\"241515\",\"passwordCredentials\":{\"password\":\"Password\",\"username\":\"chris\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, Username: "chris", Password: "Password", TenantID: "241515", Region: "region-a.geo-1"}
			return identity.Authenticate(params)
		})
}

func TestAuthenticateWithOnlyUsernameAndPasswordAndTenantNameValid(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"tenantName\":\"MyTenant\",\"passwordCredentials\":{\"password\":\"Password\",\"username\":\"chris\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, Username: "chris", Password: "Password", TenantName: "MyTenant", Region: "region-a.geo-1"}
			return identity.Authenticate(params)
		})
}

// If both TenantName and TenantID are written then the auth request will fail.
func TestAuthenticateWithTokenTenantNameAndTenantIDOnlyTenantIDTokenUsedWritten(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"tenantId\":\"MyTenantID\",\"token\":{\"id\":\"token\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, AuthToken: "token", TenantName: "MyTenant", TenantID: "MyTenantID", Region: "region-a.geo-1"}
			return identity.Authenticate(params)
		})
}

func TestAuthenticateTokenIDValid(t *testing.T) {
	verifyRequestPayloadCreated(t, "{\"auth\":{\"tenantId\":\"5678\",\"token\":{\"id\":\"token153\"}}}",
		func(url string) identity.Version2Authenticator {
			params := common.AuthenticationParameters{AuthURL: url, TenantID: "5678", AuthToken: "token153", Region: "region-a.geo-1"}
			return identity.Authenticate(params)
		})
}

func TestComputeV2GetServiceURLNoError(t *testing.T) {
	url := getServiceURLValid(t, "compute", "2", getSampleAuthPayload(validExpiringTime))
	testutil.Equals(t, "https://foo.a.compute/v2/10394455779270", url)
}

func TestComputeGetServiceURLNoMatchShouldError(t *testing.T) {
	getServiceURLErrorCondition(t, "region-c.geo-1", "compute", "1.1", getSampleAuthPayload(validExpiringTime),
		"Found serviceType 'compute' in the ServiceCatalog but cannot find an endpoint with the specified region 'region-c.geo-1' and version '1.1'")
}

func TestNetworkV2GetServiceURLNoError(t *testing.T) {
	url := getServiceURLValid(t, "network", "2", getSampleAuthPayload(validExpiringTime))
	testutil.Equals(t, "https://foo.a.neutron/v2", url)
}

func TestGetServiceURLForNonNotImplementedServiceTypeShouldError(t *testing.T) {
	getServiceURLErrorCondition(t, "region-a.geo-1", "otherService", "v2.0", getSampleAuthPayload(validExpiringTime),
		"GetServiceURL is not supported for service type 'otherService'")
}

func TestGetServiceURLNoRegionMatchShouldError(t *testing.T) {
	getServiceURLErrorCondition(t, "region-c.geo-1", "compute", "v2.0", getSampleAuthPayload(validExpiringTime),
		"Found serviceType 'compute' in the ServiceCatalog but cannot find an endpoint with the specified region 'region-c.geo-1' and version 'v2.0'")
}

func getServiceURLValid(t *testing.T, serviceType, version, authScInputPayload string) string {
	apiServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testutil.Equals(t, "POST", r.Method)
			testutil.Assert(t, strings.HasSuffix(r.URL.String(), "/tokens"), "Expected /tokens to be at the end of the url")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(authScInputPayload))

			return
		}))
	defer apiServer.Close()
	params := common.AuthenticationParameters{AuthURL: apiServer.URL, Username: "chris", Password: "Password", Region: "region-a.geo-1"}
	authenticator := identity.Authenticate(params)

	url, err := authenticator.GetServiceURL(serviceType, version)
	testutil.IsNil(t, err)
	return url
}

func getServiceURLErrorCondition(t *testing.T, region string, serviceType string, version string, authScInputPayload, expectedError string) {
	apiServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testutil.Equals(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(authScInputPayload))

			return
		}))

	defer apiServer.Close()
	params := common.AuthenticationParameters{AuthURL: apiServer.URL, Username: "chris", Password: "Password", Region: region}
	authenticator := identity.Authenticate(params)

	url, err := authenticator.GetServiceURL(serviceType, version)
	testutil.Equals(t, "", url)
	testutil.Equals(t, expectedError, err.Error())
}

func verifyRequestPayloadCreated(t *testing.T, expectedRequestPayload string, authFunc func(url string) identity.Version2Authenticator) {
	apiServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testutil.Equals(t, "POST", r.Method)
			responseByteArr, _ := ioutil.ReadAll(r.Body)
			actualRequestPayload := string(responseByteArr)
			testutil.Equals(t, expectedRequestPayload, actualRequestPayload)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(getSampleAuthPayload(validExpiringTime)))
			return
		}))

	defer apiServer.Close()
	authenticator := authFunc(apiServer.URL)
	token, err := authenticator.GetToken()
	testutil.IsNil(t, err)
	testutil.Equals(t, "HPAuth10_3b08341242f74692c29ea936f27e7a4b74", token)
}

func testValidAuth(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testutil.Equals(t, "POST", r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(getSampleAuthPayload(validExpiringTime)))
			return
		}))
}

var validExpiringTime = "2099-01-07T07:22:52.184Z"

func getSampleAuthPayload(expiringTime string) string {
	return fmt.Sprintf(sampleAuthPayloadSimplifiedServiceCatalog, expiringTime)
}

var sampleAuthPayloadSimplifiedServiceCatalog = `{
   "access":{
      "token":{
         "expires":"%s",
         "id":"HPAuth10_3b08341242f74692c29ea936f27e7a4b74",
         "tenant":{
            "id":"10394455779270",
            "name":"Hewlett-Packard3172"
         }
      },
      "user":{
         "id":"10615932248942",
         "name":"chrisrob1111",
         "otherAttributes":{
            "domainStatus":"enabled",
            "domainStatusCode":"00"
         },
         "roles":[
            {
               "id":"10419409370304",
               "serviceId":"170",
               "name":"net-admin",
               "tenantId":"10394458779270"
            }
         ]
      },
      "serviceCatalog":[
         {
            "name":"nova",
            "type":"compute",
            "endpoints":[
               {
                  "tenantId":"10394455779270",
                  "publicURL":"https:\/\/foo.a.compute\/v2\/10394455779270",
                  "region":"region-a.geo-1",
                  "versionId":"2",
                  "versionInfo":"https:\/\/foo.a.compute/v2\/",
                  "versionList":"https:\/\/foo.a.compute"
               },
               {
                  "tenantId":"10394455779270",
                  "publicURL":"https:\/\/foo.b.compute/v2\/10394455779270",
                  "region":"region-b.geo-1",
                  "versionId":"2",
                  "versionInfo":"https:\/\/foo.b.compute\/v2\/",
                  "versionList":"https:\/\/foo.b.compute"
               }
            ]
         },
		{
            "name":"nova",
            "type":"computev3",
            "endpoints":[
               {
                  "tenantId":"10394455779270",
                  "publicURL":"https:\/\/foo.a.compute\/v3\/10394455779270",
                  "region":"region-a.geo-1",
                  "versionId":"2",
                  "versionInfo":"https:\/\/foo.a.compute\/v3\/",
                  "versionList":"https:\/\/foo.a.compute"
               },
               {
                  "tenantId":"10394455779270",
                  "publicURL":"https:\/\/foo.b.compute\/v3\/10394455779270",
                  "region":"region-b.geo-1",
                  "versionId":"2",
                  "versionInfo":"https:\/\/foo.b.compute\/v3\/",
                  "versionList":"https:\/\/foo.b.compute"
               }
            ]
         },
		{
            "name":"neutron",
            "type":"network",
            "endpoints":[
               {
                  "publicURL":"https:\/\/foo.a.neutron",
                  "region":"region-a.geo-1",
                  "versionId":"3",
                  "versionInfo":"https:\/\/foo.a.neutron",
                  "versionList":"https:\/\/foo.a.neutron"
               },
               {
                  "publicURL":"https:\/\/foo.b.neutron",
                  "region":"region-b.geo-1",
                  "versionId":"3",
                  "versionInfo":"https:\/\/foo.b.neutron",
                  "versionList":"https:\/\/foo.b.neutron"
               }
            ]
         },
		{
            "name":"glance",
            "type":"image",
            "endpoints":[
               {
                  "publicURL":"http://127.0.0.1:61158/a/v1",
                  "region":"region-a.geo-1",
                  "versionInfo":"http://127.0.0.1:61158/a",
                  "versionList":"http://127.0.0.1:61158/a"
               },
               {
                  "publicURL":"http://127.0.0.1:61158/b/v1",
                  "region":"region-b.geo-1",
                  "versionInfo":"http://127.0.0.1:61158/b",
                  "versionList":"http://127.0.0.1:61158/b"
               }
            ]
         }
      ]
   }
}`

func getSamplePrivateHelionAuthPayload(expiringTime string) string {
	return fmt.Sprintf(privateHelionAuthSimplifiedServiceCatalog, expiringTime)
}

var privateHelionAuthSimplifiedServiceCatalog = `{
   "access":{
      "token":{
         "issued_at":"2015-01-09T23:00:38.658750",
         "expires":"%s",
         "id":"{SHA1}f46fe2900926256684550a6540916cf6f7ec53f1",
         "tenant":{
            "enabled":true,
            "description":"",
            "name":"devdev",
            "id":"f5d081724ab641242489e4461f776ec86"
         }
      },
      "serviceCatalog":[
         {
            "endpoints_links":[

            ],
            "endpoints":[
               {
                  "adminURL":"http://10.8.50.144:8774/v2/f5d081724ab64134b89e4461f776ec86",
                  "region":"regionOne",
                  "publicURL":"http://10.8.50.144:8774/v2/f5d081724ab64134b89e4461f776ec86",
                  "internalURL":"http://10.8.50.144:8774/v2/f5d081724ab64134b89e4461f776ec86",
                  "id":"9a483beb29864ddd9b42d7d0e1204d0e"
               }
            ],
            "type":"compute",
            "name":"nova"
         }
      ],
      "user":{
         "username":"cf-dev",
         "roles_links":[

         ],
         "id":"b41cc258f26d4f0389dfb0d4b3475b9e",
         "roles":[
            {
               "name":"_member_"
            }
         ],
         "name":"cf-dev"
      },
      "metadata":{
         "is_admin":0,
         "roles":[
            "9fe2ff9ee4384b1894a90878d3e92bab"
         ]
      }
   }
}`
