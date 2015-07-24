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

package serviceVersionList

import (
	"git.openstack.org/stackforge/golang-client.git/testUtil"

	"testing"
)

var tokn = "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"

func TestEndpoints(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, hpHelionImageVersionPayload, "")
	defer apiServer.Close()

	versions, err := Endpoints(apiServer.URL, tokn, nil)
	testUtil.IsNil(t, err)
	numVersions := len(versions)
	testUtil.Equals(t, 2, numVersions)

	testUtil.Equals(t, "v1.1", versions[0].ID)
	testUtil.Equals(t, "CURRENT", versions[0].Status)

	testUtil.Equals(t, "v1.0", versions[1].ID)
	testUtil.Equals(t, "SUPPORTED", versions[1].Status)
}

func TestFindEndpointVersionValid(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, hpHelionImageVersionPayload, "")
	defer apiServer.Close()

	url, err := FindEndpointVersion(apiServer.URL, tokn, nil, "v1.0")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "http://region-a.geo-1.images.hpcloudsvc.com/v1/", url)
}

func TestFindEndpointVersionNoValueFoundNoErrorShouldOccur(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, hpHelionImageVersionPayload, "")
	defer apiServer.Close()

	url, err := FindEndpointVersion(apiServer.URL, tokn, nil, "v2.0")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "", url)
}

func TestGetSelfLinkFound(t *testing.T) {
	v := Version{Status: "CURRENT",
		ID:    "ID",
		Links: []Link{Link{HRef: "http://loc", Rel: "self"}, Link{HRef: "http://doc", Rel: "decribedby"}}}
	foundLink := v.GetSelfLink()
	testUtil.Equals(t, foundLink.HRef, "http://loc")
}

func TestGetSelfLinkNotFoundNoErrorOccurs(t *testing.T) {
	v := Version{Status: "CURRENT",
		ID:    "ID",
		Links: []Link{Link{HRef: "http://doc", Rel: "decribedby"}}}
	foundLink := v.GetSelfLink()
	testUtil.Equals(t, foundLink.HRef, "")
}

// Actual payload from hp helion
var hpHelionImageVersionPayload = `{
   "versions":[
      {
         "status":"CURRENT",
         "id":"v1.1",
         "links":[
            {
               "href":"https://region-a.geo-1.images.hpcloudsvc.com/v1/",
               "rel":"self"
            }
         ]
      },
      {
         "status":"SUPPORTED",
         "id":"v1.0",
         "links":[
            {
               "href":"https://region-a.geo-1.images.hpcloudsvc.com/v1/",
               "rel":"self"
            }
         ]
      }
   ]
}`
