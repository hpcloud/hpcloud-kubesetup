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

package compute_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestLimits(t *testing.T) {
	expectedLimits := limitsContainer{}
	json.Unmarshal([]byte(sampleLimits), &expectedLimits)
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleLimits, "/limits")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	limits, err := service.Limits()
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, expectedLimits.Limits, limits)
}

var testTime, _ = time.Parse(`"2006-01-02T15:04:05"`, `"2012-11-27T17:24:52"`)

func TestErrorIsRateLimitShouldSucceed(t *testing.T) {

	sampleRateLimit := compute.RateLimit{
		NextAvailable: testTime,
		Remaining:     120,
		Unit:          "MINUTE",
		Value:         120,
		Verb:          "POST",
	}

	sampleBytes, _ := json.Marshal(sampleRateLimit)
	samplejson := `{ "limit": ` + string(sampleBytes) + " }"

	val := map[string][]string{"Content-Type": []string{"application/json"}}
	reader := testUtil.NewTestReadCloser([]byte(samplejson))
	response := http.Response{
		Header:     val,
		StatusCode: 413,
		Body:       reader,
	}
	httpStatus := misc.NewHTTPStatus(&response, "message")
	isRateLimit, rateLimit := compute.ErrorIsRateLimit(httpStatus)
	testUtil.Equals(t, true, isRateLimit)
	testUtil.Equals(t, sampleRateLimit, rateLimit)
}

type limitsContainer struct {
	Limits compute.Limits `json:"limits"`
}

// example derived from http://developer.openstack.org/api-ref-compute-v2-ext.html
// Using sample because of datetime values in rate limits.
var sampleLimits = `{
    "limits": {
        "absolute": {
            "maxImageMeta": 128,
            "maxPersonality": 5,
            "maxPersonalitySize": 10240,
            "maxSecurityGroupRules": 20,
            "maxSecurityGroups": 10,
            "maxServerMeta": 128,
            "maxTotalCores": 20,
            "maxTotalFloatingIps": 10,
            "maxTotalInstances": 10,
            "maxTotalKeypairs": 100,
            "maxTotalRAMSize": 51200,
            "totalCoresUsed": 0,
            "totalInstancesUsed": 0,
            "totalRAMUsed": 0,
            "totalSecurityGroupsUsed": 0,
            "totalFloatingIpsUsed": 0
        },
        "rate": [
            {
                "limit": [
                    {
                        "next-available": "2012-11-27T17:24:52Z",
                        "remaining": 120,
                        "unit": "MINUTE",
                        "value": 120,
                        "verb": "POST"
                    },
                    {
                        "next-available": "2012-11-27T17:24:52Z",
                        "remaining": 120,
                        "unit": "MINUTE",
                        "value": 120,
                        "verb": "PUT"
                    },
                    {
                        "next-available": "2012-11-27T17:24:52Z",
                        "remaining": 120,
                        "unit": "MINUTE",
                        "value": 120,
                        "verb": "DELETE"
                    }
                ],
                "regex": ".*",
                "uri": "*"
            }
        ]
    }
}`
