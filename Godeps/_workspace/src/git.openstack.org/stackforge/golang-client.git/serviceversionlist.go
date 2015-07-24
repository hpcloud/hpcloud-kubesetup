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

/*
Package serviceVersionList implements a client library for accessing service version endpoints
This can be used to understand what versions of particular services are running
*/
package serviceVersionList

import (
	"fmt"
	"net/url"

	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
)

// Version is a structure for one image version
type Version struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	Links  []Link `json:"links"`
}

// GetSelfLink returns the self link for the Version.
// this Link contains the url of the service itself.
func (version Version) GetSelfLink() Link {
	for _, link := range version.Links {
		if link.Rel == "self" {
			return link
		}
	}

	return Link{}
}

// Link contains information of a link
type Link struct {
	HRef string `json:"href"`
	Rel  string `json:"rel"`
}

// Endpoints gets the service versions from the specified endpoint.
func Endpoints(serviceURL string, tokenID string, r requester.SendRequestFunction) ([]Version, error) {
	var versions versions
	err := misc.GetJSONWithTokenAndRequester(serviceURL, tokenID, r, &versions)
	return versions.Versions, err
}

// FindEndpointVersion will get the version endpoint and seach for a current or supported version of the service.
func FindEndpointVersion(serviceURL string, tokenID string, r requester.SendRequestFunction, version string) (string, error) {
	endPointVersions, err := Endpoints(serviceURL, tokenID, r)

	if err != nil {
		return "", fmt.Errorf("Error attempting to get the service version list at url '%s' :%s", serviceURL, err)
	}

	foundEndpoint := FilterVersion(endPointVersions, "CURRENT", version)
	if foundEndpoint.ID != "" {
		return preferSchemeFromServiceCatalog(serviceURL, foundEndpoint.GetSelfLink().HRef)
	}

	foundEndpoint = FilterVersion(endPointVersions, "SUPPORTED", version)
	if foundEndpoint.ID != "" {
		return preferSchemeFromServiceCatalog(serviceURL, foundEndpoint.GetSelfLink().HRef)
	}

	return "", nil
}

// FilterVersion will locate the version of the service
func FilterVersion(versions []Version, status string, serviceVersion string) Version {
	for _, version := range versions {
		if version.Status == status && version.ID == serviceVersion {
			return version
		}
	}

	return Version{}
}

func preferSchemeFromServiceCatalog(serviceCatalogPublicURL string, serviceVersionListPublicURL string) (string, error) {
	catalogParsedURL, err := url.Parse(serviceCatalogPublicURL)
	if err != nil {
		return "", err
	}
	serviceParsedURL, err := url.Parse(serviceVersionListPublicURL)
	if err != nil {
		return "", err
	}

	if catalogParsedURL.Scheme != serviceParsedURL.Scheme {
		serviceParsedURL.Scheme = catalogParsedURL.Scheme
	}

	finalURL := serviceParsedURL.String()

	return finalURL, nil
}

// Versions is a structure for all service verisons
type versions struct {
	Versions []Version `json:"versions"`
}
