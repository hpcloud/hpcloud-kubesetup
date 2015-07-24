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

// Package testUtil has helpers to be used with unit tests
package testUtil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/misc"
)

// InvalidJSONPayload indicates an invalid payload
const InvalidJSONPayload = -1

// InvalidJSONPayloadString indicates an invalid payload
const InvalidJSONPayloadString = "{Invalid JSON payload"

const red = "red"
const reset = "reset"

var colors = map[string]string{red: "\033[31m", reset: "\033[0m"}

// Equals fails the test if exp is not equal to act.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s%s:%d:\n\n\texp: %#v\n\n\tgot: %#v%s\n\n", colors[red], filepath.Base(file), line, exp, act, colors[reset])
		tb.FailNow()
	}
}

// Assert fails the test if the condition is false.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s%s:%d: "+msg+"%s\n\n", colors[red], append([]interface{}{filepath.Base(file), line}, v...), colors[reset])
		tb.FailNow()
	}
}

// IsNil ensures that the act interface is nil
// otherwise an error is raised.
func IsNil(tb testing.TB, act interface{}) {
	if act != nil {
		tb.Error("expected nil", act)
		tb.FailNow()
	}
}

// CreateGetJSONTestRequestServer creates a httptest.Server that can be used to test GetJson requests. Just specify the token,
// json payload that is to be read by the response, and a verification func that can be used
// to do additional validation of the request that is built
func CreateGetJSONTestRequestServer(t *testing.T, expectedAuthTokenValue string, jsonResponsePayload string, verifyRequest func(*http.Request)) *httptest.Server {
	return CreateGetJSONTestRequestServerWithStatus(t, expectedAuthTokenValue, http.StatusOK, jsonResponsePayload, verifyRequest)
}

// CreateGetJSONTestRequestServerWithStatus creates a httptest.Server that can be used to test GetJson requests. Just specify
// the token, the response status, json payload that is to be read by the response, and a verification func that can be used
// to do additional validation of the request that is built
func CreateGetJSONTestRequestServerWithStatus(t *testing.T, expectedAuthTokenValue string, responseStatus int, jsonResponsePayload string, verifyRequest func(*http.Request)) *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if responseStatus == http.StatusOK {
				HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
				HeaderValuesEqual(t, r, "Accept", "application/json")
			}

			verifyRequest(r)

			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				if responseStatus == http.StatusOK {
					w.Write([]byte(jsonResponsePayload))
				}
				w.WriteHeader(responseStatus)
				return
			}

			t.Error(errors.New("Failed: r.Method == GET"))
		}))
}

// CreateGetTestRequestServer creates a httptest.Server that can be used with Get. The user is required to specify the entire body and headers of
// the response.
func CreateGetTestRequestServer(t *testing.T, expectedAuthTokenValue string, responseStatus int, headers http.Header, body []byte, verifyRequest func(*http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if responseStatus == http.StatusOK {
				HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
				HeaderValuesEqual(t, r, "Accept", "*/*")
			}

			verifyRequest(r)

			if r.Method == "GET" {
				for k, vArray := range headers {
					for _, value := range vArray {
						w.Header().Add(k, value)
					}
				}

				if responseStatus == http.StatusOK {
					w.Write(body)
				}

				w.WriteHeader(responseStatus)
				return
			}

			t.Error(errors.New("Failed: r.Method == GET"))
		}))
}

// CreateHeadTestRequestServer creates a httptest.Server that can be used with Head. The user is required to specify the entire body and headers of
// the response.
func CreateHeadTestRequestServer(t *testing.T, expectedAuthTokenValue string, responseStatus int, headers http.Header, body []byte, verifyRequest func(*http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if responseStatus == http.StatusOK {
				HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
				HeaderValuesEqual(t, r, "Accept", "*/*")
			}

			verifyRequest(r)

			if r.Method == "HEAD" {
				for k, vArray := range headers {
					for _, value := range vArray {
						w.Header().Add(k, value)
					}
				}

				if responseStatus == http.StatusOK {
					w.Write(body)
				}

				w.WriteHeader(responseStatus)
				return
			}

			t.Error(errors.New("Failed: r.Method == HEAD"))
		}))
}

// CreateGetJSONTestRequestServerWithMockObject creates a http.Server that can be used to test GetJSON requests.
// Specify the token, response object which will be marshaled to a json payload and the expected ending url.
func CreateGetJSONTestRequestServerWithMockObject(t *testing.T, token string, mockResponseObject interface{}, urlEndsWith string) *httptest.Server {
	return CreateGetJSONTestRequestServerWithMockObjectAndStatus(t, token, http.StatusOK, mockResponseObject, urlEndsWith)
}

// CreateGetJSONTestRequestServerWithMockObjectAndStatus creates a http.Server that can be used to test GetJSON requests.
// Specify the token, the response status, response object which will be marshaled to a json payload and the expected ending url.
func CreateGetJSONTestRequestServerWithMockObjectAndStatus(t *testing.T, token string, responseStatus int, mockResponseObject interface{}, urlEndsWith string) *httptest.Server {

	var mockResponse []byte
	var err error

	if responseStatus == InvalidJSONPayload {
		mockResponse = []byte(InvalidJSONPayloadString)
	} else {
		mockResponse, err = json.Marshal(mockResponseObject)
		if err != nil {
			t.Error("Test failed to marshal mockResponseObject:", err)
		}
	}

	return CreateGetJSONTestRequestServerVerifyStatusAndURL(t, token, responseStatus, string(mockResponse), urlEndsWith)
}

// CreateGetJSONTestRequestServerVerifyURL creates a http.Server that can be used to test GetJSON requests.
// Specify the token, the marshaled json payload of response object, and the expected ending url.
func CreateGetJSONTestRequestServerVerifyURL(t *testing.T, token string, jsonValue string, urlEndsWith string) *httptest.Server {
	return CreateGetJSONTestRequestServerVerifyStatusAndURL(t, token, http.StatusOK, jsonValue, urlEndsWith)
}

// CreateGetJSONTestRequestServerVerifyStatusAndURL creates a http.Server that can be used to test GetJSON requests.
// Specify the token, the response status, the marshaled json payload of response object, and the expected ending url.
func CreateGetJSONTestRequestServerVerifyStatusAndURL(t *testing.T, token string, responseStatus int, jsonValue string, urlEndsWith string) *httptest.Server {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, urlEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + urlEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	return CreateGetJSONTestRequestServerWithStatus(t, token, responseStatus, jsonValue, anon)
}

// CreatePostJSONTestRequestServer creates a http.Server that can be used to test PostJson requests. Specify the token,
// response json payload and the url and request body that is expected.
func CreatePostJSONTestRequestServer(t *testing.T, expectedAuthTokenValue string, outputResponseJSONPayload string, expectedRequestURLEndsWith string, expectedRequestBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			HeaderValuesEqual(t, r, "Accept", "application/json")
			HeaderValuesEqual(t, r, "Content-Type", "application/json")
			reqURL := r.URL.String()
			if !strings.HasSuffix(reqURL, expectedRequestURLEndsWith) {
				t.Error(errors.New("Incorrect url created, expected:" + expectedRequestURLEndsWith + " at the end, actual url:" + reqURL))
			}
			actualRequestBody := dumpRequestBody(r)
			if actualRequestBody != expectedRequestBody {
				t.Error(errors.New("Incorrect payload created, expected:'" + expectedRequestBody + "', actual '" + actualRequestBody + "'"))
			}
			if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/json")
				// Status Code has to be written before writing the payload or
				// else it defaults to 200 OK instead.
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(outputResponseJSONPayload))
				return
			}
			t.Error(errors.New("Failed: r.Method == POST"))
		}))
}

// CreateDeleteTestRequestServer creates a http.Server that can be used to test Delete requests.
func CreateDeleteTestRequestServer(t *testing.T, expectedAuthTokenValue string, urlEndsWith string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			HeaderValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			reqURL := r.URL.String()
			if !strings.HasSuffix(reqURL, urlEndsWith) {
				t.Error(errors.New("Incorrect url created, expected '" + urlEndsWith + "' at the end, actual url:" + reqURL))
			}
			if r.Method == "DELETE" {
				w.WriteHeader(204)
				return
			}

			t.Error(errors.New("Failed: r.Method == DELETE"))
		}))
}

// HeaderValuesEqual verifies the header values are equal to the expected value.
func HeaderValuesEqual(t *testing.T, req *http.Request, name string, expectedValue string) {
	actualValue := req.Header.Get(name)
	if actualValue != expectedValue {
		t.Error(fmt.Errorf("Expected Header {Name:'%s', Value:'%s', actual value '%s'", name, expectedValue, actualValue))
	}
}

// UserHomeDir returns current user home directory including ending path separator, for example, "/Users/johndoe/" on Mac.
func UserHomeDir(t *testing.T) string {
	curUser, err := user.Current()
	if err != nil {
		t.Error(fmt.Errorf("Failed to get current user: '%s'", err.Error()))
	}
	return misc.Strcat(curUser.HomeDir, string(os.PathSeparator))
}

func dumpRequestBody(request *http.Request) (body string) {
	requestWithBody, _ := httputil.DumpRequest(request, true)
	requestWithoutBody, _ := httputil.DumpRequest(request, false)
	body = strings.Replace(string(requestWithBody), string(requestWithoutBody), "", 1)
	return
}

func init() {
	// if you're running in Jenkins, reset all the colors to an empty string.
	// The colors don't actually display as colors in Jenkins, just as escape
	// sequences.
	if os.Getenv("JENKINS_URL") != "" {
		for k := range colors {
			colors[k] = ""
		}
	}
}

// TestReadCloser is a test implementation of a io.Reader
// and io.Closer
type TestReadCloser struct {
	buffer          *bytes.Buffer
	ErrorOnCallRead bool
}

// NewTestReadCloser creates a new reader with the specified bytes
func NewTestReadCloser(data []byte) TestReadCloser {
	return TestReadCloser{buffer: bytes.NewBuffer(data)}
}

// Read will read from the buffer
func (r TestReadCloser) Read(p []byte) (n int, err error) {
	if r.ErrorOnCallRead {
		return 0, fmt.Errorf("Read should not be called!")
	}
	return r.buffer.Read(p)
}

// Close is an empty impl
func (r TestReadCloser) Close() error {
	return nil
}
