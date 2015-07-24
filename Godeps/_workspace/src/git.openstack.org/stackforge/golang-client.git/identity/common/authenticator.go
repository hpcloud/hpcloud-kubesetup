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

//Package common provides interfaces and ability to allow getting tokens and service urls
// independent of the auth version.
package common

import (
	"errors"
	"fmt"
	"os"

	"git.openstack.org/stackforge/golang-client.git/misc/requester"
)

// Authenticator gives the ability to authenticate and get information
// from the service catalog in a version independent way.
type Authenticator interface {
	// GetToken returns a string that can be used as the X-Auth-Token for openstack requests or an error.
	GetToken() (string, error)

	// GetServiceURL searches the service catalog and version endpoints if needed to
	// find a service endpoint specified. If nothing is found an error should be
	// returned.
	GetServiceURL(serviceType string, version string) (string, error)
}

// AuthenticationParameters has all the parameters that are standard env vars
// for open stack using the python client.
type AuthenticationParameters struct {
	AuthURL    string
	TenantID   string
	TenantName string
	Username   string
	Password   string
	Region     string
	CACert     string
	AuthToken  string
}

// FromEnvVars will read the authURL, TenantID, TenantName
// Username, password and region from standard openstack
// env variables and create a new set of authentication parameters.
func FromEnvVars() (AuthenticationParameters, error) {
	params := AuthenticationParameters{}
	w := newGetEnvWrapper()
	params.AuthURL = w.getRequiredEnv("OS_AUTH_URL")
	params.TenantID = os.Getenv("OS_TENANT_ID")
	params.TenantName = os.Getenv("OS_TENANT_NAME")
	params.Username = w.getOptionalEnv("OS_USERNAME")
	params.Password = w.getOptionalEnv("OS_PASSWORD")
	params.Region = w.getRequiredEnv("OS_REGION_NAME")
	params.CACert = w.getOptionalEnv("OS_CACERT")
	params.AuthToken = w.getOptionalEnv("OS_AUTH_TOKEN")

	//Check here that either AuthToken OR username-password is set
	if params.AuthToken == "" {
		if params.Username == "" || params.Password == "" {
			w.Errors = append(w.Errors, fmt.Errorf("Either AuthToken or Username-Password should be provided"))
		}
	}
	return params, w.consolidateErrors()
}

type getEnvWrapper struct {
	Errors []error
}

func newGetEnvWrapper() getEnvWrapper {
	return getEnvWrapper{Errors: []error{}}
}

func (w *getEnvWrapper) getRequiredEnv(envVarName string) string {
	value := os.Getenv(envVarName)
	if value == "" {
		w.Errors = append(w.Errors, fmt.Errorf("No value provided for env variable '%s'", envVarName))
	}

	return value
}

func (w *getEnvWrapper) getOptionalEnv(envVarName string) string {
	value := os.Getenv(envVarName)

	return value
}

func (w *getEnvWrapper) appendError(error error) {
	w.Errors = append(w.Errors, error)
}

// ConsolidateErrors consolidates all errors to one error.
func (w *getEnvWrapper) consolidateErrors() error {
	consolidatedError := ""
	for _, innerError := range w.Errors {
		if consolidatedError == "" {
			consolidatedError = innerError.Error()
		} else {
			consolidatedError = consolidatedError + "\n" + innerError.Error()
		}
	}

	if len(w.Errors) > 0 {
		return errors.New(consolidatedError)
	}

	return nil
}

// SimpleAuthenticator is an authenticator for scenarios where
// hardcoded token and service url are used. Useful for testing
// scenarios.
type SimpleAuthenticator struct {
	Token           string
	ServiceURL      string
	Error           error
	requestFunction requester.SendRequestFunction
}

// NewSimpleAuthenticator creates a new authenticator with no errors.
func NewSimpleAuthenticator(token, serviceURL string) SimpleAuthenticator {
	return SimpleAuthenticator{Token: token, ServiceURL: serviceURL}
}

// GetToken returns the Token string and error in the struct.
func (a SimpleAuthenticator) GetToken() (string, error) {
	return a.Token, a.Error
}

// GetServiceURL returns the Token string and error in the struct.
func (a SimpleAuthenticator) GetServiceURL(serviceType string, version string) (string, error) {
	return a.ServiceURL, a.Error
}

// SetFunction will set the Function used to make requests.
func (a *SimpleAuthenticator) SetFunction(requestFunction requester.SendRequestFunction) {
	a.requestFunction = requestFunction
}

// Function will retrieve a function that can make a request.
func (a *SimpleAuthenticator) Function() requester.SendRequestFunction {
	return a.requestFunction
}
