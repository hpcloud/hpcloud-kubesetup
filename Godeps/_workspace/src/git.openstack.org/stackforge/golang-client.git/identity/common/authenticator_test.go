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

// authenticator_test.go
package common_test

import (
	"os"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestFromFromEnvVarsShouldErrorWithNoAuthURLUserNameAndPassword(t *testing.T) {
	savedParams, _ := common.FromEnvVars()
	if savedParams.AuthURL != "" {
		t.Skip("Skipping testing as open stack environments are set.")
	}

	os.Setenv("OS_AUTH_URL", "")
	os.Setenv("OS_TENANT_ID", "")
	os.Setenv("OS_TENANT_NAME", "")
	os.Setenv("OS_USERNAME", "")
	os.Setenv("OS_PASSWORD", "")
	os.Setenv("OS_REGION_NAME", "")

	_, err := common.FromEnvVars()
	testUtil.Assert(t, err != nil, "Expected an error")
	testUtil.Equals(t, "No value provided for env variable 'OS_AUTH_URL'\nNo value provided for env variable 'OS_REGION_NAME'\nEither AuthToken or Username-Password should be provided", err.Error())

	applyParameter(savedParams)
}

func applyParameter(p common.AuthenticationParameters) {
	os.Setenv("OS_AUTH_URL", p.AuthURL)
	os.Setenv("OS_TENANT_ID", p.TenantID)
	os.Setenv("OS_TENANT_NAME", p.TenantName)
	os.Setenv("OS_USERNAME", p.Username)
	os.Setenv("OS_PASSWORD", p.Password)
	os.Setenv("OS_REGION_NAME", p.Region)
}
