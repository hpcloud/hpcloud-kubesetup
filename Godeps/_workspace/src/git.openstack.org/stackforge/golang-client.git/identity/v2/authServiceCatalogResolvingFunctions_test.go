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
package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var testEndpoints = []endpoint{endpoint{PublicURL: "http://endpoint1", Region: "region1"},
	endpoint{PublicURL: "http://endpoint2", Region: "region2"}}

func TestServiceURLComputeV2Returns(t *testing.T) {

	sc := []service{createService("compute", "nova", testEndpoints...)}

	url, err := getComputeServiceURL("t", nil, sc, "compute", "region2", "2")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://endpoint2", url)
}

func TestServiceURLComputeV3Returns(t *testing.T) {
	sc := []service{createService("computev3", "nova", testEndpoints...)}

	url, err := getComputeServiceURL("t", nil, sc, "compute", "region1", "3")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://endpoint1", url)
}

func TestServiceURLComputeErrorsWhenNonExistant(t *testing.T) {
	url, err := getComputeServiceURL("t", nil, []service{}, "compute", "region1", "3")
	testUtil.Equals(t, "", url)
	testUtil.Assert(t, err != nil, "Expected an error")
	testUtil.Equals(t, "ServiceCatalog does not contain serviceType 'computev3'", err.Error())
}

func TestServiceURLNetworkV2Returns(t *testing.T) {
	sc := []service{createService("network", "neutron", testEndpoints...)}

	url, err := getAppendVersionServiceURL("t", nil, sc, "network", "region1", "2")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://endpoint1/v2", url)
}

func TestServiceURLNetworkV3Returns(t *testing.T) {
	sc := []service{createService("network", "neutron", testEndpoints...)}

	url, err := getAppendVersionServiceURL("t", nil, sc, "network", "region2", "3")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://endpoint2/v3", url)
}

func TestServiceURLNetworkErrorsWhenNonExistant(t *testing.T) {
	url, err := getAppendVersionServiceURL("t", nil, []service{}, "network", "region1", "3")
	testUtil.Equals(t, "", url)
	testUtil.Assert(t, err != nil, "Expected an error")
	testUtil.Equals(t, "ServiceCatalog does not contain serviceType 'network'", err.Error())
}

func TestDefaultServiceURLFoundInCatalogFoundInVersionList(t *testing.T) {
	apiServer := testVersionList(t)
	defer apiServer.Close()
	sc := []service{createService("image", "glance", endpoint{Region: "region1", PublicURL: apiServer.URL + "/publicurl", VersionList: apiServer.URL + "/versionlist"})}
	url, err := defaultGetVersionURLFilterByVersion("t", nil, sc, "image", "region1", "1")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://region-b.geo-1.image.hpcloudsvc.com/v1/", url)
}

func TestDefaultServiceURLFoundInCatalogErrorsNotFoundInVersionList(t *testing.T) {
	apiServer := testVersionList(t)
	defer apiServer.Close()
	sc := []service{createService("image", "glance", endpoint{Region: "region1", PublicURL: apiServer.URL + "/publicurl", VersionList: apiServer.URL + "/versionlist"})}
	url, err := defaultGetVersionURLFilterByVersion("t", nil, sc, "image", "region1", "2")
	testUtil.Equals(t, "", url)
	testUtil.Assert(t, err != nil, "Expected an error")
	testUtil.Equals(t, "Found serviceType 'image' in the ServiceCatalog but cannot find an endpoint with the specified region 'region1' and version '2'", err.Error())
}

func TestDefaultServiceResolverURLNotFoundInServiceCatalogError(t *testing.T) {
	apiServer := testVersionList(t)
	defer apiServer.Close()
	url, err := defaultGetVersionURLFilterByVersion("t", nil, []service{}, "image", "region1", "1")
	testUtil.Equals(t, "", url)
	testUtil.Assert(t, err != nil, "Expected an error")
	testUtil.Equals(t, "ServiceCatalog does not contain serviceType 'image'", err.Error())
}

func testVersionList(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			testUtil.Equals(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(sampleImageVersionList))

			return
		}))
}

func createService(serviceType, serviceName string, endpoints ...endpoint) service {
	return service{Name: serviceName, Type: serviceType, Endpoints: endpoints}
}

var sampleImageVersionList = `{
   "versions":[
      {
         "status":"CURRENT",
         "id":"v1.1",
         "links":[
            {
               "href":"http://region-b.geo-1.image.hpcloudsvc.com/v1/",
               "rel":"self"
            }
         ]
      },
      {
         "status":"SUPPORTED",
         "id":"v1.0",
         "links":[
            {
               "href":"http://region-b.geo-1.image.hpcloudsvc.com/v1/",
               "rel":"self"
            }
         ]
      }
   ]
}`
