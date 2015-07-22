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

// PortResponse returns a set of values of the a port response.
type PortResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	AdminStateUp        bool      `json:"admin_state_up"`
	PortSecurityEnabled bool      `json:"port_security_enabled"`
	DeviceID            string    `json:"device_id"`
	DeviceOwner         string    `json:"device_owner"`
	NetworkID           string    `json:"network_id"`
	TenantID            string    `json:"tenant_id"`
	MacAddress          string    `json:"mac_address"`
	FixedIPs            []FixedIP `json:"fixed_ips"`
	SecurityGroups      []string  `json:"security_groups"`
}

// CreatePortParameters holds a set of values that specify how
// to create a new port.
type CreatePortParameters struct {
	AdminStateUp bool      `json:"admin_state_up"`
	Name         string    `json:"name"`
	NetworkID    string    `json:"network_id"`
	FixedIPs     []FixedIP `json:"fixed_ips"`
}

// PortResponses is a type for a slice of PortResponses.
type PortResponses []PortResponse

// FixedIP is holds data that specifies a fixed IP.
type FixedIP struct {
	SubnetID  string `json:"subnet_id,omitempty"`
	IPAddress string `json:"ip_address"`
}

// Ports issues a GET request that returns the found port responses
func (networkService Service) Ports() ([]PortResponse, error) {
	p := portsResp{}
	reqURL, err := networkService.buildRequestURL("/ports")
	if err != nil {
		return p.Ports, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &p)
	if err != nil {
		return nil, err
	}

	return p.Ports, nil
}

// Port issues a GET request that returns a specific port response.
func (networkService Service) Port(id string) (PortResponse, error) {
	port := portResp{}
	reqURL, err := networkService.buildRequestURL("/ports/", id)
	if err != nil {
		return port.Port, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &port)
	return port.Port, err
}

// DeletePort issues a DELETE to the specified port url to delete it.
func (networkService Service) DeletePort(id string) error {
	reqURL, err := networkService.buildRequestURL("/ports/", id)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, networkService.authenticator)
}

// CreatePort issues a POST to create the specified port and return a PortResponse.
func (networkService Service) CreatePort(parameters CreatePortParameters) (PortResponse, error) {
	parametersContainer := createPortContainer{Port: parameters}
	portResponse := portResp{}
	reqURL, err := networkService.buildRequestURL("/ports")
	if err != nil {
		return portResponse.Port, err
	}

	err = misc.PostJSON(reqURL, networkService.authenticator, parametersContainer, &portResponse)
	return portResponse.Port, err
}

type portsResp struct {
	Ports []PortResponse `json:"ports"`
}

type portResp struct {
	Port PortResponse `json:"port"`
}

type createPortContainer struct {
	Port CreatePortParameters `json:"port"`
}
