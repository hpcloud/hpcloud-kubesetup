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
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// FloatingIP contains properties is of an ip
// that can be associated with a server and can
// be in a specified pool.
type FloatingIP struct {
	ID         string `json:"id"`
	IP         string `json:"ip"`
	InstanceID string `json:"instance_id"`
	FixedIP    string `json:"fixed_ip"`
	Pool       string `json:"pool"`
}

// FloatingIPs will issue a get query that returns a list of floating ips.
func (computeService Service) FloatingIPs() ([]FloatingIP, error) {
	var r = floatingIPsResp{}
	reqURL, err := computeService.buildRequestURL("/os-floating-ips")
	if err != nil {
		return r.FloatingIPs, err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.FloatingIPs, err
}

// FloatingIP will issue a get query that returns a floating ip.
func (computeService Service) FloatingIP(id string) (FloatingIP, error) {
	var r = floatingIPResp{}
	reqURL, err := computeService.buildRequestURL("/os-floating-ips/", id)
	if err != nil {
		return r.FloatingIP, err
	}
	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.FloatingIP, err
}

// CreateFloatingIP will issue a Post query creates a new floating ip and return the created value.
func (computeService Service) CreateFloatingIP(pool string) (FloatingIP, error) {
	var requestParameters = createFloatingIPRequest{Pool: pool}
	var r = floatingIPResp{}
	reqURL, err := computeService.buildRequestURL("/os-floating-ips")
	if err != nil {
		return r.FloatingIP, err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, requestParameters, &r)

	return r.FloatingIP, err
}

// CreateFloatingIPInDefaultPool will issue a Post query to create a floating ip in the default pool
// and returns the created value.
func (computeService Service) CreateFloatingIPInDefaultPool() (FloatingIP, error) {
	return computeService.CreateFloatingIP("")
}

// DeleteFloatingIP will issue a delete query to delete the floating ip.
func (computeService Service) DeleteFloatingIP(floatingIPID string) (err error) {
	reqURL, err := computeService.buildRequestURL("/os-floating-ips/", floatingIPID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

type floatingIPResp struct {
	FloatingIP FloatingIP `json:"floating_ip"`
}

type floatingIPsResp struct {
	FloatingIPs []FloatingIP `json:"floating_ips,omitempty"`
}

type createFloatingIPRequest struct {
	Pool string `json:"pool"`
}
