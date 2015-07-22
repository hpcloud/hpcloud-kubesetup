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

package network_test

import (
	"testing"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
	network "git.openstack.org/stackforge/golang-client.git/network/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

func TestSecurityGroupScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	authParameters, err := common.FromEnvVars()
	if err != nil {
		t.Fatal("There was an error getting authentication env variables:", err)
	}

	authenticator := identity.Authenticate(authParameters)
	authenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))

	networkService := network.NewService(&authenticator)

	sampleSecurityGroup := network.CreateSecurityGroupParameters{
		Description: "Description",
		Name:        "NewCustomSecurityGroup",
	}

	securityGroup, err := networkService.CreateSecurityGroup(sampleSecurityGroup)
	if err != nil {
		t.Fatal("Cannot create security group:", err)
	}

	tcp := network.TCP
	cidr := "0.0.0.0/0"
	sampleRule := network.CreateSecurityGroupRuleParameters{
		PortRangeMin:    intPointer(80),
		PortRangeMax:    intPointer(80),
		IPProtocol:      &tcp,
		RemoteIPPrefix:  &cidr,
		SecurityGroupID: securityGroup.ID,
		Direction:       "ingress",
	}

	securityGroupRule, err := networkService.CreateSecurityGroupRule(sampleRule)
	if err != nil {
		t.Error("Cannot create security group rule:", err)
	}

	testUtil.Equals(t, 80, *securityGroupRule.PortRangeMin)
	testUtil.Equals(t, 80, *securityGroupRule.PortRangeMax)
	testUtil.Equals(t, cidr, *securityGroupRule.RemoteIPPrefix)

	sampleRule2 := network.CreateSecurityGroupRuleParameters{
		PortRangeMin:    nil,
		PortRangeMax:    nil,
		IPProtocol:      &tcp,
		SecurityGroupID: securityGroup.ID,
		Direction:       "ingress",
		RemoteGroupID:   &securityGroup.ID,
	}

	securityGroupRule2, err := networkService.CreateSecurityGroupRule(sampleRule2)
	if err != nil {
		t.Error("Cannot create security group rule:", err)
	}

	testUtil.Equals(t, (*int)(nil), securityGroupRule2.PortRangeMin)
	testUtil.Equals(t, (*int)(nil), securityGroupRule2.PortRangeMax)
	testUtil.Equals(t, securityGroup.ID, *securityGroupRule2.RemoteGroupID)
	testUtil.Equals(t, (*string)(nil), securityGroupRule2.RemoteIPPrefix)

	queriedSecurityGroup, err := networkService.SecurityGroup(securityGroup.ID)
	if err != nil {
		t.Error("Cannot requery security group:", err)
	}

	testUtil.Equals(t, 4, len(queriedSecurityGroup.SecurityGroupRules))

	err = networkService.DeleteSecurityGroupRule(securityGroupRule.ID)
	if err != nil {
		t.Error("Cannot delete security group rule:", err)
	}

	err = networkService.DeleteSecurityGroup(queriedSecurityGroup.ID)
	if err != nil {
		t.Error("Cannot delete security group:", err)
	}
}

func intPointer(v int) *int {
	return &v
}
