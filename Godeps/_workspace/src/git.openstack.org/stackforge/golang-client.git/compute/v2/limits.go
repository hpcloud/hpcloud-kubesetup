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
	"encoding/json"
	"time"

	"git.openstack.org/stackforge/golang-client.git/misc"
)

// AbsoluteLimits is a structure for all properties of for the limits
type AbsoluteLimits struct {
	MaxImageMeta            int `json:"maxImageMeta"`
	MaxPersonality          int `json:"maxPersonality"`
	MaxPersonalitySize      int `json:"maxPersonalitySize"`
	MaxSecurityGroupRules   int `json:"maxSecurityGroupRules"`
	MaxSecurityGroups       int `json:"maxSecurityGroups"`
	MaxServerMeta           int `json:"maxServerMeta"`
	MaxTotalCores           int `json:"maxTotalCores"`
	MaxTotalFloatingIps     int `json:"maxTotalFloatingIps"`
	MaxTotalInstances       int `json:"maxTotalInstances"`
	MaxTotalKeypairs        int `json:"maxTotalKeypairs"`
	MaxTotalRAMSize         int `json:"maxTotalRAMSize"`
	TotalCoresUsed          int `json:"totalCoresUsed"`
	TotalInstancesUsed      int `json:"totalInstancesUsed"`
	TotalRAMUsed            int `json:"totalRAMUsed"`
	TotalSecurityGroupsUsed int `json:"totalSecurityGroupsUsed"`
	TotalFloatingIpsUsed    int `json:"totalFloatingIpsUsed"`
}

// Limits has the absolute and Rate limits for the tenant
type Limits struct {
	AbsoluteLimits AbsoluteLimits        `json:"absolute"`
	Rate           []RateLimitsContainer `json:"rate"`
}

// RateLimitsContainer has the rate limits.
type RateLimitsContainer struct {
	Limit []RateLimit `json:"limit"`
	Regex string      `json:"regex"`
	URL   string      `json:"uri"`
}

// RateLimit has properties for the rate limit.
type RateLimit struct {
	NextAvailable time.Time `json:"next-available"`
	Remaining     int       `json:"remaining"`
	Unit          string    `json:"unit"`
	Value         int       `json:"value"`
	Verb          string    `json:"verb"`
}

// ErrorIsRateLimit will probe the error to see if its a RateLimit
func ErrorIsRateLimit(err error) (bool, RateLimit) {
	httpStatus, ok := err.(misc.HTTPStatus)
	if ok && httpStatus.StatusCode == 413 && misc.ContentTypeIsJSON(httpStatus.Header) {
		body, err := httpStatus.GetBody()
		if err == nil {
			c := singleRateLimitContainer{}
			err := json.Unmarshal(body, &c)
			if err == nil {
				return true, c.RateLimit
			}
		}
	}

	return false, RateLimit{}
}

// Limits will issue a get request to retrieve the all limits for the user.
func (computeService Service) Limits() (Limits, error) {
	var limits = limitsContainer{}
	reqURL, err := computeService.buildRequestURL("/limits")
	if err != nil {
		return limits.Limits, err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &limits)
	return limits.Limits, err
}

type limitsContainer struct {
	Limits Limits `json:"limits"`
}

type singleRateLimitContainer struct {
	RateLimit RateLimit `json:"limit"`
}
