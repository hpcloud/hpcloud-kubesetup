// Copyright (c) 2015 Hewlett-Packard Development Company, L.P.
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

// Package network is used to create, delete, and query, networks, ports and subnets
package network

import "git.openstack.org/stackforge/golang-client.git/misc"

// IPProtocol is the type that can be
// ICMP, TCP, or UDP
type IPProtocol string

const (
	// ICMP is a constant for IPProtocol type.
	ICMP IPProtocol = "icmp"
	// TCP is a constant for IPProtocol type.
	TCP IPProtocol = "tcp"
	// UDP is a constant for IPProtocol type.
	UDP IPProtocol = "udp"
)

// SecurityGroup represents a Security Group in Neutron
type SecurityGroup struct {
	ID                 string              `json:"id"`
	TenantID           string              `json:"tenant_id"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`
}

// SecurityGroupRule represents a Security Group Rule in Neutron
type SecurityGroupRule struct {
	ID              string      `json:"id"`
	Direction       string      `json:"direction"`
	IPProtocol      *IPProtocol `json:"protocol"`
	EtherType       string      `json:"ethertype"`
	TenantID        string      `json:"tenant_id"`
	PortRangeMin    *int        `json:"port_range_min"`
	PortRangeMax    *int        `json:"port_range_max"`
	SecurityGroupID string      `json:"security_group_id"`
	RemoteGroupID   *string     `json:"remote_group_id"`
	RemoteIPPrefix  *string     `json:"remote_ip_prefix"`
}

// CreateSecurityGroupParameters contains the required parameters for creating a
// new security group.
type CreateSecurityGroupParameters struct {
	TenantID    string `json:"tenant_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateSecurityGroupRuleParameters represents the request to create a new security group rule.
type CreateSecurityGroupRuleParameters struct {
	Direction       string      `json:"direction"`
	PortRangeMin    *int        `json:"port_range_min"`
	PortRangeMax    *int        `json:"port_range_max"`
	IPProtocol      *IPProtocol `json:"protocol"`
	SecurityGroupID string      `json:"security_group_id"`
	RemoteGroupID   *string     `json:"remote_group_id"`
	RemoteIPPrefix  *string     `json:"remote_ip_prefix"`
}

// securityGroupsResponse represents a response from Neutron for multiple security groups
type securityGroupsResponse struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
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

// SecurityGroups will retrieve a list of security groups from Neutron
func (networkService Service) SecurityGroups() ([]SecurityGroup, error) {
	respContainer := securityGroupsResponse{}
	reqURL, err := networkService.buildRequestURL("/security-groups")
	if err != nil {
		return []SecurityGroup{}, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &respContainer)
	if err != nil {
		return []SecurityGroup{}, err
	}

	return respContainer.SecurityGroups, nil
}

// SecurityGroup will retrieve a single security group from Neutron
func (networkService Service) SecurityGroup(id string) (SecurityGroup, error) {
	respContainer := securityGroupContainer{}
	reqURL, err := networkService.buildRequestURL("/security-groups/", id)
	if err != nil {
		return SecurityGroup{}, err
	}

	err = misc.GetJSON(reqURL, networkService.authenticator, &respContainer)
	if err != nil {
		return SecurityGroup{}, err
	}

	return respContainer.SecurityGroup, nil
}

// CreateSecurityGroup will create a new security group in Neutron
func (networkService Service) CreateSecurityGroup(parameters CreateSecurityGroupParameters) (SecurityGroup, error) {
	var requestParameters = createSecurityGroupContainer{SecurityGroup: parameters}
	var r = securityGroupContainer{}
	reqURL, err := networkService.buildRequestURL("/security-groups")
	if err != nil {
		return r.SecurityGroup, err
	}

	err = misc.PostJSON(reqURL, networkService.authenticator, requestParameters, &r)

	return r.SecurityGroup, err
}

// DeleteSecurityGroup will issue a delete query to delete the specified security group.
func (networkService Service) DeleteSecurityGroup(securityGroupID string) (err error) {
	reqURL, err := networkService.buildRequestURL("/security-groups/", securityGroupID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, networkService.authenticator)
}

// CreateSecurityGroupRule creates a new security group rule in the specified parent.
func (networkService Service) CreateSecurityGroupRule(sgr CreateSecurityGroupRuleParameters) (SecurityGroupRule, error) {
	ruleContainer := securityGroupRuleCreateContainer{SecurityGroupRule: sgr}
	r := securityGroupRuleContainer{}
	reqURL, err := networkService.buildRequestURL("/security-group-rules")
	if err != nil {
		return r.SecurityGroupRule, err
	}

	err = misc.PostJSON(reqURL, networkService.authenticator, ruleContainer, &r)

	return r.SecurityGroupRule, err
}

// DeleteSecurityGroupRule will issue a delete query to delete the specified security group rule.
func (networkService Service) DeleteSecurityGroupRule(securityGroupRuleID string) (err error) {
	reqURL, err := networkService.buildRequestURL("/security-group-rules/", securityGroupRuleID)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, networkService.authenticator)
}
