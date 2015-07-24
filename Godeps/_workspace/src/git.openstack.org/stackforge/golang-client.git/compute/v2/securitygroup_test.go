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
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var icmpValue = compute.ICMP
var sampleGroup = compute.Group{Name: "name", TenantID: "tenantId"}
var sampleIPRange = compute.IPRange{CIDR: "10.3.1.1/8"}

// security group rules

var sampleCIDRRule = compute.SecurityGroupRule{
	FromPort:      133,
	Group:         nil,
	ID:            "f2a15c2-1111-49d7-ac80-01eba4641111",
	IPProtocol:    &icmpValue,
	IPRange:       sampleIPRange,
	ParentGroupID: "a226372b-3333-40c4-b4ed-a16e3cb92222",
	ToPort:        436}

var sampleGroupRule = compute.SecurityGroupRule{
	FromPort:      133,
	Group:         &sampleGroup,
	ID:            "f2a15c2-1111-49d7-ac80-01eba4641111",
	IPProtocol:    &icmpValue,
	IPRange:       sampleIPRange,
	ParentGroupID: "a226372b-3333-40c4-b4ed-a16e3cb92222",
	ToPort:        436}

var sampleCIDRRuleBytes, _ = json.Marshal(sampleCIDRRule)
var sampleCIDRRuleJSONResponse = `{ "security_group_rule":` + string(sampleCIDRRuleBytes) + `}`

var sampleGroupRuleBytes, _ = json.Marshal(sampleGroupRule)
var sampleGroupRuleJSONResponse = `{ "security_group_rule":` + string(sampleGroupRuleBytes) + `}`

// security groups

var sampleSecurityGroup = compute.SecurityGroup{
	Description: "Description",
	ID:          "a226372b-3333-40c4-b4ed-a16e3cb92222",
	Name:        "Name",
	Rules:       []compute.SecurityGroupRule{sampleCIDRRule},
	TenantID:    "tenantid"}
var sampleSecurityGroupBytes, _ = json.Marshal(sampleSecurityGroup)

var sampleSecurityGroupJSONResponse = `{ "security_group":` + string(sampleSecurityGroupBytes) + `}`
var sampleSecurityGroupsJSONResponse = `{ "security_groups": [` + string(sampleSecurityGroupBytes) + `] }`

func TestGetSecurityGroups(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleSecurityGroupsJSONResponse, "/os-security-groups")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	sgs, err := service.SecurityGroups()
	testUtil.IsNil(t, err)
	testUtil.Equals(t, len(sgs), 1)
	testUtil.Equals(t, sampleSecurityGroup, sgs[0])
}

func TestGetServerSecurityGroups(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleSecurityGroupsJSONResponse, "/os-security-groups/servers/aghagr")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	sgs, err := service.ServerSecurityGroups("aghagr")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, len(sgs), 1)
	testUtil.Equals(t, sampleSecurityGroup, sgs[0])
}

func TestGetSecurityGroup(t *testing.T) {
	apiServer := testUtil.CreateGetJSONTestRequestServerVerifyURL(t, tokn, sampleSecurityGroupJSONResponse, "/os-security-groups/w9236264")
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	sg, err := service.SecurityGroup("w9236264")
	testUtil.IsNil(t, err)
	testUtil.Equals(t, sampleSecurityGroup, sg)
}

func TestDeleteSecurityGroup(t *testing.T) {
	name := "securitygroup"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/os-security-groups/"+name)
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	err := service.DeleteSecurityGroup(name)
	testUtil.IsNil(t, err)
}

func TestCreateSecurityGroup(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleSecurityGroupJSONResponse, "/os-security-groups",
		`{"security_group":{"tenant_id":"tenantID","addSecurityGroup":"addsecgroup","name":"name","description":"Description"}}`)
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	parameters := compute.CreateSecurityGroupParameters{TenantID: "tenantID", AddSecurityGroup: "addsecgroup",
		Name: "name", Description: "Description"}

	securityGroup, err := service.CreateSecurityGroup(parameters)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, sampleSecurityGroup, securityGroup)
}

func TestDeleteSecurityGroupRule(t *testing.T) {
	name := "securitygrouprule"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/os-security-group-rules/"+name)
	defer apiServer.Close()

	service := CreateComputeService(apiServer.URL)
	err := service.DeleteSecurityGroupRule(name)
	testUtil.IsNil(t, err)
}

func TestCreateSecurityGroupRuleWithCIDR(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleCIDRRuleJSONResponse, "/os-security-group-rules",
		`{"security_group_rule":{"from_port":133,"to_port":436,"ip_protocol":"ICMP","group_id":null,"parent_group_id":"a226372b-3333-40c4-b4ed-a16e3cb92222","cidr":"10.3.1.1/8"}}`)
	defer apiServer.Close()

	cidr := "10.3.1.1/8"

	service := CreateComputeService(apiServer.URL)
	createRule := compute.CreateSecurityGroupRuleParameters{
		FromPort:      133,
		IPProtocol:    icmpValue,
		ParentGroupID: "a226372b-3333-40c4-b4ed-a16e3cb92222",
		ToPort:        436,
		CIDR:          &cidr}

	securityGroupRule, err := service.CreateSecurityGroupRule(createRule)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, sampleCIDRRule, securityGroupRule)
}

func TestCreateSecurityGroupRuleWithGroup(t *testing.T) {
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, sampleGroupRuleJSONResponse, "/os-security-group-rules",
		`{"security_group_rule":{"from_port":133,"to_port":436,"ip_protocol":"ICMP","group_id":"19ffd595-daba-49f0-91f0-741bac706a3a","parent_group_id":"a226372b-3333-40c4-b4ed-a16e3cb92222","cidr":null}}`)
	defer apiServer.Close()

	groupID := "19ffd595-daba-49f0-91f0-741bac706a3a"

	service := CreateComputeService(apiServer.URL)
	createRule := compute.CreateSecurityGroupRuleParameters{
		FromPort:      133,
		IPProtocol:    icmpValue,
		ParentGroupID: "a226372b-3333-40c4-b4ed-a16e3cb92222",
		ToPort:        436,
		GroupID:       &groupID}

	securityGroupRule, err := service.CreateSecurityGroupRule(createRule)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, sampleGroupRule, securityGroupRule)
}

func TestCompabilitySecurityGroupPayloads(t *testing.T) {
	// Read private cloud response content from file
	securityGroupsTestFilePath := "./testdata/securitygroup_test_sgs.json"
	securityGroupsTestFileContent, err := ioutil.ReadFile(securityGroupsTestFilePath)
	if err != nil {
		t.Error(fmt.Errorf("Failed to read JSON file %s: '%s'", securityGroupsTestFilePath, err.Error()))
	}

	// Decode the content
	sampleSGs := compute.SecurityGroupsContainer{}
	err = json.Unmarshal(securityGroupsTestFileContent, &sampleSGs)
	if err != nil {
		t.Error(fmt.Errorf("Failed to decode JSON file %s: '%s'", securityGroupsTestFilePath, err.Error()))
	}

	// Test the SDK API computeService.SecurityGroups()
	anon := func(computeService *compute.Service) {

		securityGroups, err := computeService.SecurityGroups()
		if err != nil {
			t.Error(err)
		}

		if len(securityGroups) != len(sampleSGs.SecurityGroups) {
			t.Error(errors.New("Incorrect number of securityGroups found"))
		}

		// Verify returned availability zones match original sample availability zones
		testUtil.Equals(t, sampleSGs.SecurityGroups, securityGroups)
	}

	testComputeServiceAction(t, "os-security-groups", string(securityGroupsTestFileContent), anon)
}
