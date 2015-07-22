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

// SecurityGroup is a structure has properties of
// the security group from open stack.
type SecurityGroup struct {
	Description string              `json:"description,omitempty"`
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Rules       []SecurityGroupRule `json:"rules,omitempty"`
	TenantID    string              `json:"tenant_id,omitempty"`
}

// CreateSecurityGroupParameters contains
// the required parameters for creating a
// new security group.
type CreateSecurityGroupParameters struct {
	TenantID         string `json:"tenant_id"`
	AddSecurityGroup string `json:"addSecurityGroup"`
	Name             string `json:"name"`
	Description      string `json:"description"`
}

// CreateSecurityGroupRuleParameters represents the request to create a new security group rule.
type CreateSecurityGroupRuleParameters struct {
	FromPort      int32      `json:"from_port"`
	ToPort        int32      `json:"to_port"`
	IPProtocol    IPProtocol `json:"ip_protocol"`
	GroupID       *string    `json:"group_id"`
	ParentGroupID string     `json:"parent_group_id"`
	CIDR          *string    `json:"cidr"`
}

// SecurityGroupRule contains the parameters
// for specifying which network is allowed
// for the rule.
type SecurityGroupRule struct {
	FromPort      int32       `json:"from_port"`
	Group         *Group      `json:"group,omitempty"`
	ID            string      `json:"id,omitempty"`
	IPProtocol    *IPProtocol `json:"ip_protocol"`
	IPRange       IPRange     `json:"ip_range,omitempty"`
	ParentGroupID string      `json:"parent_group_id"`
	ToPort        int         `json:"to_port"`
}

// SecurityGroupsContainer contains a list of security groups. Used to response to
// json as a container wrapper.
type SecurityGroupsContainer struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

// Group contains properties of the name
// and TenantID of the Security Group
type Group struct {
	TenantID string `json:"tenant_id"`
	Name     string `json:"from_port"`
}

// IPProtocol is the type that can be
// ICMP, TCP, or UDP
type IPProtocol string

const (
	// ICMP is a constant for IPProtocol type.
	ICMP IPProtocol = "ICMP"
	// TCP is a constant for IPProtocol type.
	TCP IPProtocol = "TCP"
	// UDP is a constant for IPProtocol type.
	UDP IPProtocol = "UDP"
)

// IPRange contains the CIDR
// information.
type IPRange struct {
	CIDR string `json:"cidr"`
}

// SecurityGroups will issue a get query that returns a list of security groups in the system.
func (computeService Service) SecurityGroups() ([]SecurityGroup, error) {
	var r = SecurityGroupsContainer{}
	reqURL, err := computeService.buildRequestURL("/os-security-groups")
	if err != nil {
		return r.SecurityGroups, err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.SecurityGroups, err
}

// SecurityGroup will issue a get query that returns a security group based on the id.
func (computeService Service) SecurityGroup(id string) (SecurityGroup, error) {
	var r = securityGroupContainer{}
	reqURL, err := computeService.buildRequestURL("/os-security-groups/", id)
	if err != nil {
		return r.SecurityGroup, err
	}
	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.SecurityGroup, err
}

// ServerSecurityGroups will issue a get query that returns the security groups for the specified server.
func (computeService Service) ServerSecurityGroups(serverID string) ([]SecurityGroup, error) {
	var r = SecurityGroupsContainer{}
	reqURL, err := computeService.buildRequestURL("/os-security-groups", "/servers/", serverID)
	if err != nil {
		return r.SecurityGroups, err
	}
	err = misc.GetJSON(reqURL, computeService.authenticator, &r)

	return r.SecurityGroups, err
}

// CreateSecurityGroup will issue a Post query creates a new security group and returns the created value.
func (computeService Service) CreateSecurityGroup(parameters CreateSecurityGroupParameters) (SecurityGroup, error) {
	var requestParameters = createSecurityGroupContainer{SecurityGroup: parameters}
	var r = securityGroupContainer{}
	reqURL, err := computeService.buildRequestURL("/os-security-groups")
	if err != nil {
		return r.SecurityGroup, err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, requestParameters, &r)

	return r.SecurityGroup, err
}

// DeleteSecurityGroup will issue a delete query to delete the specified security group.
func (computeService Service) DeleteSecurityGroup(securityGroupID string) (err error) {
	reqURL, err := computeService.buildRequestURL("/os-security-groups/", securityGroupID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

// CreateSecurityGroupRule creates a new security group rule in the specified parent.
func (computeService Service) CreateSecurityGroupRule(sgr CreateSecurityGroupRuleParameters) (SecurityGroupRule, error) {
	ruleContainer := securityGroupRuleCreateContainer{SecurityGroupRule: sgr}
	r := securityGroupRuleContainer{}
	reqURL, err := computeService.buildRequestURL("/os-security-group-rules")
	if err != nil {
		return r.SecurityGroupRule, err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, ruleContainer, &r)

	return r.SecurityGroupRule, err
}

// DeleteSecurityGroupRule will issue a delete query to delete the specified security group rule.
func (computeService Service) DeleteSecurityGroupRule(securityGroupRuleID string) (err error) {
	reqURL, err := computeService.buildRequestURL("/os-security-group-rules/", securityGroupRuleID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

type securityGroupContainer struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

type createSecurityGroupContainer struct {
	SecurityGroup CreateSecurityGroupParameters `json:"security_group"`
}

type securityGroupRuleContainer struct {
	SecurityGroupRule SecurityGroupRule `json:"security_group_rule"`
}

type securityGroupRuleCreateContainer struct {
	SecurityGroupRule CreateSecurityGroupRuleParameters `json:"security_group_rule"`
}
