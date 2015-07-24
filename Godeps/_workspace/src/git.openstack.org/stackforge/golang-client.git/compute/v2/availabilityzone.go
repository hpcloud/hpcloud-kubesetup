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

// AvailabilityZone contains information of OpenStack availabilityZone
type AvailabilityZone struct {
	ZoneState ZoneState `json:"zoneState"`
	ZoneName  string    `json:"zoneName"`
}

// ZoneState contains state of OpenStack availability zone
type ZoneState struct {
	Available bool `json:"available"`
}

// AvailabilityZonesContainer contains information for availability zones.
type AvailabilityZonesContainer struct {
	AvailabilityZones []AvailabilityZone `json:"availabilityZoneInfo"`
}

// AvailabilityZones will issue a GET request to retrieve all availability zones.
func (computeService Service) AvailabilityZones() ([]AvailabilityZone, error) {

	var container = AvailabilityZonesContainer{}
	url, err := computeService.buildRequestURL("/os-availability-zone")
	if err != nil {
		return container.AvailabilityZones, err
	}

	err = misc.GetJSON(url, computeService.authenticator, &container)
	if err != nil {
		return nil, err
	}

	return container.AvailabilityZones, err
}
