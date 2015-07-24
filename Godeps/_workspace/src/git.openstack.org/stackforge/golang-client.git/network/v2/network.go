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

// Package network is used to create, delete, and query, networks, ports and subnets
package network

import (
	"fmt"
	"net/url"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Service holds state that is use to make requests and get responses for networks,
// ports and subnets
type Service struct {
	authenticator common.Authenticator
}

// NewService creates a new network service client.
func NewService(authenticator common.Authenticator) Service {
	return Service{authenticator: authenticator}
}

// QueryParameters is a structure that
// contains the filter parameters for
// a network.
type QueryParameters struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	RouterExternal bool   `json:"router:external"`
	AdminStateUp   bool   `json:"admin_state_up"`
	Shared         bool   `json:"shared"`
}

func (networkService Service) buildRequestURL(suffixes ...string) (string, error) {
	serviceURL, err := networkService.authenticator.GetServiceURL("network", "2.0")
	if err != nil {
		return "", err
	}

	urlPaths := append([]string{serviceURL}, suffixes...)
	return misc.Strcat(urlPaths...), nil
}

// Response returns a set of values of the a network response.
type Response struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	Subnets             []string `json:"subnets"`
	TenantID            string   `json:"tenant_id"`
	RouterExternal      bool     `json:"router:external"`
	AdminStateUp        bool     `json:"admin_state_up"`
	Shared              bool     `json:"shared"`
	PortSecurityEnabled bool     `json:"port_security_enabled"`
}

// CreateNetworkParameters contains parameters required to create a network.
type CreateNetworkParameters struct {
	Name         string `json:"name"`
	AdminStateUp bool   `json:"admin_state_up"`
	Shared       bool   `json:"shared"`
	TenantID     string `json:"tenant_id"`
}

// EndpointLink return the HREF and Rel of an Endpoint.
type EndpointLink struct {
	HREF string `json:"href"`
	Rel  string `json:"rel"`
}

// Networks will issue a get query that returns a list of networks
func (networkService Service) Networks() ([]Response, error) {
	nwsContainer := networksResp{}
	reqURL, err := networkService.buildRequestURL("/networks")
	if err != nil {
		return nwsContainer.Networks, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &nwsContainer)
	return nwsContainer.Networks, err
}

// QueryNetworks will issue a get query that returns a list of networks
func (networkService Service) QueryNetworks(q QueryParameters) ([]Response, error) {
	nwsContainer := networksResp{}
	reqURL, err := networkService.buildNetworksQueryURL(q)
	if err != nil {
		return nwsContainer.Networks, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &nwsContainer)
	return nwsContainer.Networks, err
}

// Network will issue a get request for a specific network.
func (networkService Service) Network(id string) (Response, error) {
	nwContainer := networkResp{}
	reqURL, err := networkService.buildRequestURL("/networks/", id)
	if err != nil {
		return nwContainer.Network, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &nwContainer)
	return nwContainer.Network, err
}

// NetworkIDsByName will return the networks ID for that have the specified network name.
func (networkService Service) NetworkIDsByName(name string) ([]string, error) {
	nwContainer := networksResp{}
	reqURL, err := networkService.buildRequestURL("/networks?fields=id&name=", name)
	if err != nil {
		return []string{}, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &nwContainer)
	if err != nil {
		return []string{}, err
	}

	ids := []string{}
	for _, n := range nwContainer.Networks {
		ids = append(ids, n.ID)
	}
	return ids, err
}

// CreateNetwork will send a POST request to create a new network with the specified parameters.
func (networkService Service) CreateNetwork(parameters CreateNetworkParameters) (Response, error) {
	nContainer := networkResp{}
	reqURL, err := networkService.buildRequestURL("/networks")
	if err != nil {
		return nContainer.Network, err
	}
	err = misc.PostJSON(reqURL, networkService.authenticator, createNetworkValuesContainer{Network: parameters}, &nContainer)
	return nContainer.Network, err
}

// DeleteNetwork will delete the specified network.
func (networkService Service) DeleteNetwork(name string) (err error) {
	reqURL, err := networkService.buildRequestURL("/networks/", name)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, networkService.authenticator)
}

func (networkService Service) buildNetworksQueryURL(q QueryParameters) (string, error) {
	// Parse to create a URL structure which query parameters and the path will be encoded onto it.
	// Usage of this ensures correct encoding of url strings.
	serviceURL, err := networkService.buildRequestURL("/networks")
	if err != nil {
		return "", err
	}

	reqURL, err := url.Parse(serviceURL)
	if err != nil {
		return "", err
	}
	values := url.Values{}
	if q.AdminStateUp {
		values.Set("admin_state_up", fmt.Sprintf("%t", q.AdminStateUp))
	}
	if q.Name != "" {
		values.Set("name", q.Name)
	}
	if q.RouterExternal {
		values.Set("router:external", fmt.Sprintf("%t", q.RouterExternal))
	}
	if q.Shared {
		values.Set("shared", fmt.Sprintf("%t", q.Shared))
	}
	if q.Status != "" {
		values.Set("status", q.Status)
	}

	if len(values) > 0 {
		reqURL.RawQuery = values.Encode()
	}

	return reqURL.String(), nil
}

type createNetworkValuesContainer struct {
	Network CreateNetworkParameters `json:"network"`
}

type networksResp struct {
	Networks []Response `json:"networks"`
}

type networkResp struct {
	Network Response `json:"network"`
}
