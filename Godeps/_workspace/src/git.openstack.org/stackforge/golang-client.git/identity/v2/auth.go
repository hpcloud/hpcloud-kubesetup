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

//Package identity provides functions for client-side access to OpenStack
//IdentityService.
package identity

import (
	"fmt"
	"strings"
	"sync"
	"time"

	common "git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
)

// ServiceEndpointWithSpecifiedRegionAndVersionNotFound contains the specified error message.
const ServiceEndpointWithSpecifiedRegionAndVersionNotFound = "Found serviceType '%s' in the ServiceCatalog but cannot find an endpoint with the specified region '%s' and version '%s'"

// ServiceTypeNotFound contains the specified error message.
const ServiceTypeNotFound = "ServiceCatalog does not contain serviceType '%s'"

// ServiceTypeEndpointNotFoundWithSpecifiedRegion contains the specified error message.
const ServiceTypeEndpointNotFoundWithSpecifiedRegion = "Found serviceType '%s' in the ServiceCatalog but cannot find an endpoint with the specified region '%s'"

const password = "password"
const username = "username"

// NoServiceFoundError happens when a service is requested, but it doesn't exist in the environment
type NoServiceFoundError struct {
	msg         string // description of error
	ServiceName string // the name of the service we're looking for
	Region      string // the region this error occurred in, could be ""
	Version     string // the version requested, could be ""
}

func (e NoServiceFoundError) Error() string { return e.msg }

// AuthenticateFromEnvVars will build a Authenticator using the env vars. If
// enableHTTPDebugging is true then the DebuggingHTTPRequester is enabled
// which will print out request/response pairs to Standard out.
func AuthenticateFromEnvVars() (common.Authenticator, error) {
	authParameters, err := common.FromEnvVars()
	if err != nil {
		return nil, err
	}

	authenticator := Authenticate(authParameters)

	return &authenticator, nil
}

// Authenticate will build a Authenticator using the Authentication Params.
func Authenticate(params common.AuthenticationParameters) Version2Authenticator {
	ac := authRequestContainer{}
	ac.Auth.PasswordCredentials = map[string]string{}
	authToken := params.AuthToken
	if params.TenantID == "" {
		ac.Auth.TenantName = params.TenantName
	} else {
		ac.Auth.TenantID = params.TenantID
	}

	if authToken != "" {
		//Auth Token present
		ac.Auth.Token = map[string]string{}
		ac.Auth.Token["id"] = authToken
		return newVersion2Authenticator(params.AuthURL, ac, params.Region)
	}
	//No Auth Token, use username-password
	ac.Auth.PasswordCredentials[password] = params.Password
	ac.Auth.PasswordCredentials[username] = params.Username

	return newVersion2Authenticator(params.AuthURL, ac, params.Region)

}

// Version2Authenticator implements GetToken and GetServiceURL for a identity v2.
type Version2Authenticator struct {
	region               string
	cachingAuthRequester cachingAuthRequester
}

func newVersion2Authenticator(authURL string, ac authRequestContainer, region string) Version2Authenticator {
	caRequest := newCachingAuthRequesterRequester(authURL, ac)
	return Version2Authenticator{region: region, cachingAuthRequester: caRequest}
}

// GetToken will return the token for requests.
func (version2Authenticator *Version2Authenticator) GetToken() (string, error) {
	return version2Authenticator.cachingAuthRequester.getToken()
}

// GetServiceURL will return the url for a particular version of a service.
func (version2Authenticator *Version2Authenticator) GetServiceURL(serviceType string, version string) (string, error) {
	return version2Authenticator.cachingAuthRequester.getServiceURL(serviceType, version2Authenticator.region, version)
}

// SetFunction will set the Requester used to make requests.
func (version2Authenticator *Version2Authenticator) SetFunction(requesterFunction requester.SendRequestFunction) {
	version2Authenticator.cachingAuthRequester.requestFunction = requesterFunction
}

// Function will make a request.
func (version2Authenticator *Version2Authenticator) Function() requester.SendRequestFunction {
	return version2Authenticator.cachingAuthRequester.requestFunction
}

type cachingAuthRequester struct {
	auth                        authResponse
	authError                   error
	mutex                       *sync.Mutex
	authenticated               bool
	authURL                     string
	ac                          authRequestContainer
	serviceURLCache             map[string]string
	serviceTypeResolvingFuncMap map[string]serviceURLResolvingFunction
	requestFunction             requester.SendRequestFunction
}

func newCachingAuthRequesterRequester(authURL string, ac authRequestContainer) cachingAuthRequester {
	authURL = removeEndingSlash(authURL) + "/tokens"
	return cachingAuthRequester{
		auth:                        authResponse{},
		mutex:                       &sync.Mutex{},
		authURL:                     authURL,
		ac:                          ac,
		serviceURLCache:             map[string]string{},
		serviceTypeResolvingFuncMap: createStandardServiceTypeResolvingFuncMap(),
	}
}

func (authRequester *cachingAuthRequester) executeAuthenticationRequest() (authResponse, error) {
	ar := authResponse{}
	err := misc.PostJSONWithTokenAndRequester(authRequester.authURL, "", authRequester.requestFunction, authRequester.ac, &ar)

	if err != nil {
		return ar, err
	}

	if !ar.Access.Token.Expires.After(time.Now()) {
		return ar, fmt.Errorf("Error: The auth token has an invalid expiration.")
	}

	return ar, nil
}

func (authRequester *cachingAuthRequester) authenticateIfRequired() (t string, sc []service, cache map[string]string, ae error) {
	authRequester.mutex.Lock()
	if !authRequester.authenticated || authRequester.auth.Access.Token.Expires.Before(time.Now()) {
		authRequester.auth, authRequester.authError = authRequester.executeAuthenticationRequest()
		authRequester.authenticated = true
		authRequester.serviceURLCache = map[string]string{}
	}

	t = authRequester.auth.Access.Token.ID
	ae = authRequester.authError
	sc = authRequester.auth.Access.ServiceCatalog
	cache = authRequester.serviceURLCache
	authRequester.mutex.Unlock()

	return
}

func (authRequester *cachingAuthRequester) getToken() (string, error) {
	t, _, _, ae := authRequester.authenticateIfRequired()

	return t, ae
}

func (authRequester *cachingAuthRequester) addServiceURLToCache(key, serviceURL string) {
	authRequester.mutex.Lock()
	authRequester.serviceURLCache[key] = serviceURL
	authRequester.mutex.Unlock()
}

// getServiceURL will get the ServiceURL by first looking in the cache to see if a value exists.
// If its not in the cache then a lookup is done to see if a specific service type resolving function exists.
// If the function exists its executed to resolve the ServiceURL, and it will be cached if found.
// If the service type is not resolved then return an error indicating its not supported yet.
func (authRequester *cachingAuthRequester) getServiceURL(serviceType string, region string, version string) (string, error) {
	t, sc, cache, ae := authRequester.authenticateIfRequired()
	if ae != nil {
		return "", ae
	}

	cacheKey := serviceType + region + version
	foundCachedServiceURL, ok := cache[cacheKey]
	if ok {
		return foundCachedServiceURL, nil
	}

	resolved, foundServiceURL, err := authRequester.getServiceTypeSpecificResolvingFunc(t, sc, serviceType, region, version)
	if resolved {
		authRequester.addServiceURLToCache(cacheKey, foundServiceURL)
		return foundServiceURL, err
	}

	return "", fmt.Errorf("GetServiceURL is not supported for service type '%s'", serviceType)
}

func (authRequester *cachingAuthRequester) getServiceTypeSpecificResolvingFunc(tokenID string, sc []service, serviceType, region, version string) (found bool, serviceURL string, err error) {

	resolvingFunction, ok := authRequester.serviceTypeResolvingFuncMap[serviceType]
	if ok {
		serviceURL, err = resolvingFunction(tokenID, authRequester.requestFunction, sc, serviceType, region, version)
		found = true
	}

	return
}

func removeEndingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		return strings.TrimRight(url, "/")
	}

	return url
}

type authRequestContainer struct {
	Auth authRequest `json:"auth"`
}

type authRequest struct {
	TenantName          string            `json:"tenantName,omitempty"`
	TenantID            string            `json:"tenantId,omitempty"`
	PasswordCredentials map[string]string `json:"passwordCredentials,omitempty"`
	Token               map[string]string `json:"token,omitempty"`
}

type authResponse struct {
	Access access
}

type access struct {
	Token          token
	User           user
	ServiceCatalog []service
}

type token struct {
	ID      string
	Expires time.Time
	Tenant  tenant
}

type tenant struct {
	ID   string
	Name string
}

type user struct {
	ID         string
	Name       string
	Roles      []role
	RolesLinks []string
}

type role struct {
	ID       string
	Name     string
	TenantID string
}

type service struct {
	Name           string
	Type           string
	Endpoints      []endpoint
	EndpointsLinks []string
}

type endpoint struct {
	TenantID    string
	PublicURL   string
	InternalURL string
	Region      string
	VersionID   string
	VersionInfo string
	VersionList string
}
