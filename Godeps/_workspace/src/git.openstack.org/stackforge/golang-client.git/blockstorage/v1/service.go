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

package blockstorage

import (
	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// Service is a client service that can make
// requests against a OpenStack compute v2 service.
type Service struct {
	authenticator common.Authenticator
}

// NewService creates a new compute service client.
func NewService(authenticator common.Authenticator) Service {
	return Service{authenticator: authenticator}
}

func (blockStorageService Service) serviceURL() (string, error) {
	return blockStorageService.authenticator.GetServiceURL("volume", "1")
}

func (blockStorageService Service) buildRequestURL(suffixes ...string) (string, error) {
	serviceURL, err := blockStorageService.serviceURL()
	if err != nil {
		return "", err
	}

	urlPaths := append([]string{serviceURL}, suffixes...)
	return misc.Strcat(urlPaths...), nil
}
