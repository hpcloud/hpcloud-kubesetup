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
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestSecurityGroupScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	authenticator, err := identity.AuthenticateFromEnvVars()
	if err != nil {
		t.Fatal("Cannot authenticate from env vars:", err)
	}

	computeService := compute.NewService(authenticator)

	sampleSecurityGroup := compute.CreateSecurityGroupParameters{
		Description: "Description",
		Name:        "NewCustomSecurityGroup",
	}

	securityGroup, err := computeService.CreateSecurityGroup(sampleSecurityGroup)
	if err != nil {
		t.Fatal("Cannot create security group:", err)
	}

	tcp := compute.TCP
	cidr := "0.0.0.0/0"
	sampleRule := compute.CreateSecurityGroupRuleParameters{
		FromPort:      80,
		ToPort:        80,
		IPProtocol:    tcp,
		CIDR:          &cidr,
		ParentGroupID: securityGroup.ID,
	}

	securityGroupRule, err := computeService.CreateSecurityGroupRule(sampleRule)
	if err != nil {
		t.Fatal("Cannot create security group rule:", err)
	}

	queriedSecurityGroup, err := computeService.SecurityGroup(securityGroup.ID)
	if err != nil {
		t.Fatal("Cannot requery security group:", err)
	}

	testUtil.Assert(t, len(queriedSecurityGroup.Rules) > 0, "Expected Security group to have a rule")

	err = computeService.DeleteSecurityGroupRule(securityGroupRule.ID)
	if err != nil {
		t.Fatal("Cannot delete security group rule:", err)
	}

	err = computeService.DeleteSecurityGroup(queriedSecurityGroup.ID)
	if err != nil {
		t.Fatal("Cannot delete security group:", err)
	}
}
