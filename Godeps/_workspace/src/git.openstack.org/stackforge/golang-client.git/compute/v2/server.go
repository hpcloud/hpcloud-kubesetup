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

// Package compute is used to manage hypervisors, virtual machines, and SSH key pairs
package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
)

// defaultLimit is used as the limit for calls to the Compute service,
// for pagination purposes
const defaultLimit = 21

///////////////////////////////////////////////////////////////////////////////
// Public Structs
///////////////////////////////////////////////////////////////////////////////

// Service is a client service that can make
// requests against a OpenStack compute v2 service.
type Service struct {
	authenticator common.Authenticator
}

// NewService creates a new compute service client.
func NewService(authenticator common.Authenticator) Service {
	return Service{authenticator: authenticator}
}

func (computeService Service) serviceURL() (string, error) {
	return computeService.authenticator.GetServiceURL("compute", "2")
}

func (computeService Service) buildRequestURL(suffixes ...string) (string, error) {
	serviceURL, err := computeService.serviceURL()
	if err != nil {
		return "", err
	}

	urlPaths := append([]string{serviceURL}, suffixes...)
	return misc.Strcat(urlPaths...), nil
}

// ServerNetworkParameters contains server network related parameters
type ServerNetworkParameters struct {
	UUID    string `json:"uuid,omitempty"`
	Port    string `json:"port,omitempty"`
	FixedIP string `json:"fixed_ip,omitempty"`
}

// ServerBlockDeviceMappingV2 contains properties related to booting a new
// server using devices in block storage.
type ServerBlockDeviceMappingV2 struct {
	DeviceName          string `json:"device_name"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type,omitempty"`
	DeleteOnTermination bool   `json:"delete_on_termination,string"`
	GuestFormat         string `json:"guest_format"`
	BootIndex           int    `json:"boot_index,string"`
	UUID                string `json:"uuid,omitempty"`
	ConfigDrive         bool   `json:"config_drive,omitempty"`
}

// ServerCreationParameters contains server service parameters
type ServerCreationParameters struct {
	Name        string `json:"name"`
	ImageRef    string `json:"imageRef"`
	KeyPairName string `json:"key_name"`
	FlavorRef   string `json:"flavorRef"`

	// optional parameters
	MaxCount             *int32                       `json:"maxcount,omitempty"`
	MinCount             *int32                       `json:"mincount,omitempty"`
	UserData             *string                      `json:"user_data,omitempty"`
	AvailabilityZone     *string                      `json:"availability_zone,omitempty"`
	Networks             []ServerNetworkParameters    `json:"networks,omitempty"`
	SecurityGroups       []SecurityGroup              `json:"security_groups,omitempty"`
	Metadata             map[string]string            `json:"metadata,omitempty"`
	Personality          []ServerPersonality          `json:"personality,omitempty"`
	BlockDeviceMappingV2 []ServerBlockDeviceMappingV2 `json:"block_device_mapping_v2,omitempty"`
}

// ServerPersonality contains server personaility information.
type ServerPersonality struct {
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

// CreateServerResponse contains ID, links, admin password of the newly
// created server and the associated security groups the server was
// created with.
type CreateServerResponse struct {
	ID             string          `json:"id"`
	Links          []Link          `json:"links"`
	AdminPass      string          `json:"adminPass"`
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

// Server contains basic information of a server
type Server struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

// ServerDetail contains detailed information of a server
type ServerDetail struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Status           string               `json:"status"`
	Created          time.Time            `json:"created"`
	Updated          *time.Time           `json:"updated"`
	HostID           string               `json:"hostId"`
	Addresses        map[string][]Address `json:"addresses"`
	Links            []Link               `json:"links"`
	Image            ImageWrapper         `json:"image"`
	Flavor           Flavor               `json:"flavor"`
	TaskState        string               `json:"OS-EXT-STS:task_state"`
	VMState          string               `json:"OS-EXT-STS:vm_state"`
	PowerState       int                  `json:"OS-EXT-STS:power_state"`
	AvailabilityZone string               `json:"OS-EXT-AZ:availability_zone"`
	UserID           string               `json:"user_id"`
	TenantID         string               `json:"tenant_id"`
	AccessIPv4       string               `json:"accessIPv4"`
	AccessIPv6       string               `json:"accessIPv6"`
	ConfigDrive      string               `json:"config_drive"`
	Progress         int                  `json:"progress"`
	MetaData         map[string]string    `json:"metadata"`
	AdminPass        string               `json:"adminPass"`
	KeyName          string               `json:"key_name"`
}

// ServerDetailQueryParameters contains properties
// that allow filtering and paging of servers.
type ServerDetailQueryParameters struct {
	Name   string `json:"name"`
	Limit  int64  `json:"limit"`
	Marker string `json:"marker"`
}

// QueryParameters contains properties that allow filtering and paging of compute requests.
type QueryParameters struct {
	Limit  int64  `json:"limit"`
	Marker string `json:"marker"`
}

// Link contains information of a link
type Link struct {
	HRef string `json:"href"`
	Rel  string `json:"rel"`
}

// ImageWrapper exists because an image can be returned
// as an empty string or an Image type
type ImageWrapper struct {
	Image *Image
}

// Image contains information of an image
type Image struct {
	ID    string `json:"id"`
	Links []Link `json:"links"`
}

// Address contains address information
type Address struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type"`
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
}

///////////////////////////////////////////////////////////////////////////////
// Public Functions
///////////////////////////////////////////////////////////////////////////////

// UnmarshalJSON converts the bytes give to a Image
func (r *ImageWrapper) UnmarshalJSON(data []byte) error {
	if data == nil {
		r = &ImageWrapper{}

		return nil
	}

	if data != nil && len(data) == 2 && data[0] == '"' && data[1] == '"' {
		r = &ImageWrapper{}

		return nil
	}

	r.Image = &Image{}

	err := json.NewDecoder(bytes.NewReader(data)).Decode(r.Image)

	return err
}

// MarshalJSON converts a Image to a []byte.
func (r ImageWrapper) MarshalJSON() ([]byte, error) {
	var doc bytes.Buffer

	if r.Image == nil {
		return []byte("\"\""), nil
	}

	err := json.NewEncoder(&doc).Encode(r.Image)

	if err != nil {
		return nil, err
	}

	return doc.Bytes(), nil
}

// CreateServer creates a virtual machine cloud server
func (computeService Service) CreateServer(serverCreationParameters ServerCreationParameters) (CreateServerResponse, error) {
	sc := createServerContainer{}
	reqURL, err := computeService.buildRequestURL("/servers")
	if err != nil {
		return sc.CreateServer, err
	}

	c := serverCreateParametersContainer{ServerCreationParameters: serverCreationParameters}
	err = misc.PostJSON(reqURL, computeService.authenticator, c, &sc)
	return sc.CreateServer, err
}

// DeleteServer deletes a virtual machine cloud server
func (computeService Service) DeleteServer(id string) (err error) {
	reqURL, err := computeService.buildRequestURL("/servers/", id)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

// ServerAction invokes an action on a virtual machine cloud server
func (computeService Service) ServerAction(id string, action string, key string, value string) (err error) {
	actionKeyValue := make(map[string]string)
	actionKeyValue[key] = value

	serverActionParameter := make(serverActionCreationParameter)
	serverActionParameter[action] = actionKeyValue

	reqURL, err := computeService.buildRequestURL("/servers/", id, "/action")
	if err != nil {
		return err
	}

	var resp interface{}

	return misc.PostJSON(reqURL, computeService.authenticator, serverActionParameter, &resp)
}

// Servers retrieves basic server information
func (computeService Service) Servers() ([]Server, error) {
	result := []Server{}
	marker := ""

	more := true

	var err error
	for more {
		container := serversContainer{}
		reqURL, err := computeService.buildServersDetailQueryURL(ServerDetailQueryParameters{
			Limit:  int64(defaultLimit),
			Marker: marker,
		}, "/servers")
		if err != nil {
			return container.Servers, err
		}

		err = misc.GetJSON(reqURL.String(), computeService.authenticator, &container)

		if err != nil {
			return result, err
		}

		if len(container.Servers) < defaultLimit {
			more = false
		}

		for _, server := range container.Servers {
			result = append(result, server)
			marker = server.ID
		}
	}

	return result, err
}

// ServerDetails retrieves detailed information of servers
func (computeService Service) ServerDetails() (serverDetails []ServerDetail, err error) {
	return computeService.QueryServersDetail(ServerDetailQueryParameters{})
}

// QueryServersDetail retrieves detailed information of servers based on the query parameters provided
func (computeService Service) QueryServersDetail(queryParameters ServerDetailQueryParameters) ([]ServerDetail, error) {
	result := []ServerDetail{}
	marker := ""

	more := true

	var err error
	for more {
		params := queryParameters
		params.Limit = defaultLimit
		params.Marker = marker

		container := serversDetailContainer{}
		reqURL, err := computeService.buildServersDetailQueryURL(params, "/servers/detail")
		if err != nil {
			return nil, err
		}

		err = misc.GetJSON(reqURL.String(), computeService.authenticator, &container)

		if err != nil {
			return result, err
		}

		if len(container.ServersDetail) < defaultLimit {
			more = false
		}

		for _, server := range container.ServersDetail {
			result = append(result, server)
			marker = server.ID
		}
	}

	return result, err
}

// ServerDetail retrieves detailed information of a server
func (computeService Service) ServerDetail(id string) (serverDetail ServerDetail, err error) {
	c := serverDetailContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", id)
	if err != nil {
		return c.ServerDetail, err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &c)
	serverDetail = c.ServerDetail
	return
}

// ServerMetadata retrieves all metadata properties of the specified server.
func (computeService Service) ServerMetadata(serverID string) (map[string]string, error) {
	m := serverMetadataContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/metadata")
	if err != nil {
		return m.Metadata, err
	}
	err = misc.GetJSON(reqURL, computeService.authenticator, &m)
	return m.Metadata, err
}

// ServerMetadataItem retrieves the specific metadata item of the specified server.
func (computeService Service) ServerMetadataItem(serverID string, metaItemName string) (string, error) {
	m := serverMetaItemContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/metadata/", metaItemName)
	if err != nil {
		return "", err
	}

	err = misc.GetJSON(reqURL, computeService.authenticator, &m)
	if err != nil {
		return "", err
	}

	return m.Meta[metaItemName], nil
}

// SetServerMetadataItem sets the specific metadata item to the a value for the specified server.
// It returns the value that was specified.
func (computeService Service) SetServerMetadataItem(serverID string, metaItemName string, value string) (map[string]string, error) {
	input := serverMetaItemContainer{Meta: make(map[string]string)}
	input.Meta[metaItemName] = value
	output := serverMetaItemContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/metadata/", metaItemName)
	if err != nil {
		return make(map[string]string, 0), err
	}

	err = misc.PutJSON(reqURL, computeService.authenticator, input, &output)
	if err != nil {
		return make(map[string]string, 0), err
	}

	return output.Meta, nil
}

// SetServerMetadata sets the metadata to the specified server.
func (computeService Service) SetServerMetadata(serverID string, metadata map[string]string) (map[string]string, error) {
	input := serverMetadataContainer{Metadata: metadata}
	output := serverMetadataContainer{}
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/metadata")
	if err != nil {
		return make(map[string]string, 0), err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, input, &output)
	if err != nil {
		return make(map[string]string, 0), err
	}

	return output.Metadata, nil
}

// DeleteServerMetadataItem deletes the specified metaItem from the metadata the server.
func (computeService Service) DeleteServerMetadataItem(serverID string, metaItemName string) error {
	reqURL, err := computeService.buildRequestURL("/servers/", serverID, "/metadata/", metaItemName)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

// DeleteServerMetadata deletes the metadata of the server by sending an empty metadata to the specified server.
func (computeService Service) DeleteServerMetadata(serverID string) error {
	_, err := computeService.SetServerMetadata(serverID, make(map[string]string, 0))
	return err
}

func (computeService Service) buildServersDetailQueryURL(queryParameters ServerDetailQueryParameters, partialURLPath string) (*url.URL, error) {
	// Parse to create a URL structure which query parameters and the path will be encoded onto it.
	// Usage of this ensures correct encoding of url strings.
	serviceURL, err := computeService.serviceURL()
	if err != nil {
		return nil, err
	}

	reqURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	if queryParameters.Name != "" {
		values.Set("name", queryParameters.Name)
	}
	if queryParameters.Limit != 0 {
		values.Set("limit", fmt.Sprintf("%d", queryParameters.Limit))
	}
	if queryParameters.Marker != "" {
		values.Set("marker", queryParameters.Marker)
	}

	if len(values) > 0 {
		reqURL.RawQuery = values.Encode()
	}

	reqURL.Path += partialURLPath

	return reqURL, nil
}

func (computeService Service) buildPaginatedQueryURL(queryParameters QueryParameters, partialURLPath string) (*url.URL, error) {
	// Parse to create a URL structure which query parameters and the path will be encoded onto it.
	// Usage of this ensures correct encoding of url strings.
	serviceURL, err := computeService.serviceURL()
	if err != nil {
		return nil, err
	}

	reqURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	if queryParameters.Limit != 0 {
		values.Set("limit", fmt.Sprintf("%d", queryParameters.Limit))
	}
	if queryParameters.Marker != "" {
		values.Set("marker", queryParameters.Marker)
	}

	if len(values) > 0 {
		reqURL.RawQuery = values.Encode()
	}

	reqURL.Path += partialURLPath

	return reqURL, nil
}

///////////////////////////////////////////////////////////////////////////////
// Private Structs
///////////////////////////////////////////////////////////////////////////////

type serverCreateParametersContainer struct {
	ServerCreationParameters ServerCreationParameters `json:"server"`
}

type serversContainer struct {
	Servers []Server `json:"servers"`
}

type serversDetailContainer struct {
	ServersDetail []ServerDetail `json:"servers"`
}

type serverDetailContainer struct {
	ServerDetail ServerDetail `json:"server"`
}

type createServerContainer struct {
	CreateServer CreateServerResponse `json:"server"`
}

// serverActionCreationParameter defines the parameter for creating server action
type serverActionCreationParameter map[string]map[string]string

type serverMetadataContainer struct {
	Metadata map[string]string `json:"metadata"`
}

type serverMetaItemContainer struct {
	Meta map[string]string `json:"meta"`
}
