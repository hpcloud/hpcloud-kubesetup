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

package identity

import (
	"fmt"
	"strings"

	core "git.openstack.org/stackforge/golang-client.git"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
)

type serviceURLResolvingFunction func(tokenID string, r requester.SendRequestFunction, sc []service, serviceType, region, version string) (string, error)

func createStandardServiceTypeResolvingFuncMap() map[string]serviceURLResolvingFunction {
	return map[string]serviceURLResolvingFunction{
		"compute":      getComputeServiceURL,
		"network":      getAppendVersionServiceURL,
		"object-store": getPublicURL,
		"database":     getPublicURL,
		"image":        defaultGetVersionURLFilterByVersion,
		"volume":       getPublicURL,
	}
}

func defaultGetVersionURLFilterByVersion(tokenID string, r requester.SendRequestFunction, sc []service, serviceType, region, version string) (string, error) {
	foundService, err := findService(sc, serviceType)
	if err != nil {
		return "", err
	}

	foundVersionListEndpoints := []string{}
	for _, endPointInstance := range foundService.Endpoints {
		if endPointInstance.Region != region {
			continue
		}

		if endPointInstance.VersionList == "" {
			foundVersionListEndpoints = append(foundVersionListEndpoints, endPointInstance.PublicURL)
		} else {
			foundVersionListEndpoints = append(foundVersionListEndpoints, endPointInstance.VersionList)
		}
	}

	numEps := len(foundVersionListEndpoints)

	if numEps == 0 {
		return "", fmt.Errorf(ServiceTypeEndpointNotFoundWithSpecifiedRegion, serviceType, region)
	}

	for _, versionListURL := range foundVersionListEndpoints {
		foundServiceURL, err := core.FindEndpointVersion(versionListURL, tokenID, r, fixupVersionForEndpoints(version))
		if err != nil {
			return "", err
		}

		if foundServiceURL != "" {
			return foundServiceURL, nil
		}
	}

	return "", &NoServiceFoundError{
		msg:         fmt.Sprintf(ServiceEndpointWithSpecifiedRegionAndVersionNotFound, serviceType, region, version),
		Region:      region,
		ServiceName: serviceType,
		Version:     version}
}

func getAppendVersionServiceURL(tokenID string, r requester.SendRequestFunction, sc []service, serviceType, region, version string) (string, error) {
	foundService, err := findService(sc, serviceType)
	if err != nil {
		return "", err
	}

	for _, endPointInstance := range foundService.Endpoints {
		if endPointInstance.Region == region {
			return removeEndingSlash(endPointInstance.PublicURL) + "/v" + version, nil
		}
	}

	return "", &NoServiceFoundError{
		msg:         fmt.Sprintf(ServiceEndpointWithSpecifiedRegionAndVersionNotFound, serviceType, region, version),
		Region:      region,
		ServiceName: serviceType,
		Version:     version}
}

func getComputeServiceURL(tokenID string, r requester.SendRequestFunction, sc []service, serviceType, region, version string) (string, error) {
	if serviceType == "compute" && version == "3" {
		// NOTE from Python nova client about this oddness:
		// (cyeoh): Having the service type dependent on the API version
		// is pretty ugly, but we have to do this because traditionally the
		// catalog entry for compute points directly to the V2 API rather than
		// the root, and then doing version discovery.
		serviceType = "computev3"
	}

	return getPublicURL(tokenID, r, sc, serviceType, region, version)
}

func getPublicURL(tokenID string, r requester.SendRequestFunction, sc []service, serviceType, region, version string) (string, error) {
	foundService, err := findService(sc, serviceType)
	if err != nil {
		return "", err
	}

	for _, endPointInstance := range foundService.Endpoints {
		if endPointInstance.Region == region {
			return endPointInstance.PublicURL, nil
		}
	}

	return "", &NoServiceFoundError{
		msg:         fmt.Sprintf(ServiceEndpointWithSpecifiedRegionAndVersionNotFound, serviceType, region, version),
		Region:      region,
		ServiceName: serviceType,
		Version:     version}
}

func findService(sc []service, serviceType string) (service, error) {
	foundService := service{}
	for _, v := range sc {
		if v.Type == serviceType {
			return v, nil
		}
	}

	return foundService, &NoServiceFoundError{
		msg:         fmt.Sprintf(ServiceTypeNotFound, serviceType),
		ServiceName: serviceType}
}

// Whenever selecting a version from the endpoint list
// the version the version needs to be formatted a particular
// way for the string match to work. This normalizes the
// version strings so that matches will occur.
func fixupVersionForEndpoints(version string) string {
	updatedVersion := version
	if !strings.HasPrefix(version, "v") {
		updatedVersion = "v" + updatedVersion
	}
	if !strings.Contains(version, ".") {
		updatedVersion = updatedVersion + ".0"
	}

	return updatedVersion
}
