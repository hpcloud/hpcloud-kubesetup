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

package network_test

import (
	"encoding/json"
	"errors"
	"testing"

	network "git.openstack.org/stackforge/golang-client.git/network/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var sampleExternalGateway = network.ExternalGatewayInfo{
	NetworkID:  "networkid",
	EnableSNAT: true,
}

var sampleRouter = network.Router{
	ExternalGatewayInfo: sampleExternalGateway,
	Status:              "ACTIVE",
	AdminStateUp:        true,
	TenantID:            "tenantID",
	Name:                "RouterName",
	ID:                  "ID",
}

var sampleRouterJsonBytes, _ = json.Marshal(sampleRouter)
var sampleRouterJsonInRoutersContainer = `{ "routers": [` + string(sampleRouterJsonBytes) + `]}`
var sampleRouterJsonInRouterContainer = `{ "router":` + string(sampleRouterJsonBytes) + `}`

func TestGetRouters(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleRouterJsonInRoutersContainer, "/routers")
	defer apiServer.Close()

	service := CreateNetworkService(apiServer.URL)
	routers, err := service.Routers()
	if err != nil {
		t.Error(err)
	}

	if len(routers) != 1 {
		t.Error(errors.New("Error: Expected 2 keypairs to be listed"))
	}
	testUtil.Equals(t, sampleRouter, routers[0])
}

func TestGetRouter(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleRouterJsonInRouterContainer, "/routers/RID")
	defer apiServer.Close()

	service := CreateNetworkService(apiServer.URL)
	router, err := service.Router("RID")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleRouter, router)
}

func TestDeleteRouter(t *testing.T) {
	id := "RID"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/routers/"+id)
	defer apiServer.Close()

	service := CreateNetworkService(apiServer.URL)
	err := service.DeleteRouter(id)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateRouter(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleRouterJsonInRouterContainer, "/routers",
		`{"router":{"external_gateway_info":{"network_id":"networkid"}}}`)
	defer apiServer.Close()

	service := CreateNetworkService(apiServer.URL)
	actualRouter, err := service.CreateRouter("networkid")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleRouter, actualRouter)
}
