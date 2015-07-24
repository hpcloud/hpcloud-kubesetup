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

package network

import (
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Router is a structure has properties of
// the router from open stack.
type Router struct {
	Status              string              `json:"status"`
	ExternalGatewayInfo ExternalGatewayInfo `json:"external_gateway_info"`
	Name                string              `json:"name"`
	AdminStateUp        bool                `json:"admin_state_up"`
	TenantID            string              `json:"tenant_id"`
	ID                  string              `json:"id"`
}

// ExternalGatewayInfo indicates a router is an external gateway.
type ExternalGatewayInfo struct {
	NetworkID  string `json:"network_id"`
	EnableSNAT bool   `json:"enable_snat,omitempty"`
}

// Routers will issue a get query that returns a list of security groups in the system.
func (networkService Service) Routers() ([]Router, error) {
	r := routersContainer{}
	reqURL, err := networkService.buildRequestURL("/routers")
	if err != nil {
		return r.Routers, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &r)

	return r.Routers, err
}

// Router will issue a get query that returns a security group based on the id.
func (networkService Service) Router(id string) (Router, error) {
	r := routerContainer{}
	reqURL, err := networkService.buildRequestURL("/routers/", id)
	if err != nil {
		return r.Router, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &r)

	return r.Router, err
}

// CreateRouter will issue a Post query creates a new security group and returns the created value.
func (networkService Service) CreateRouter(networkID string) (Router, error) {
	var requestParameters = createRouterContainer{CreateRouter: createRouter{GatewayInfo: ExternalGatewayInfo{NetworkID: networkID}}}
	var r = routerContainer{}

	reqURL, err := networkService.buildRequestURL("/routers")
	if err != nil {
		return r.Router, err
	}

	err = misc.PostJSON(reqURL, networkService.authenticator, requestParameters, &r)

	return r.Router, err
}

// DeleteRouter will issue a delete query to delete the specified security group.
func (networkService Service) DeleteRouter(routerID string) (err error) {
	reqURL, err := networkService.buildRequestURL("/routers/", routerID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, networkService.authenticator)
}

type routersContainer struct {
	Routers []Router `json:"routers"`
}

type routerContainer struct {
	Router Router `json:"router"`
}

type createRouterContainer struct {
	CreateRouter createRouter `json:"router"`
}

type createRouter struct {
	GatewayInfo ExternalGatewayInfo `json:"external_gateway_info"`
}
