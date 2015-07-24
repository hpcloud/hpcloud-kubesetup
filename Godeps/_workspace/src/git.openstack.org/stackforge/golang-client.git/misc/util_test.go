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

package misc_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var token = "2350971-5716-8165"
var authenticator = common.SimpleAuthenticator{Token: token}

func TestDelete(t *testing.T) {
	var apiServer = testUtil.CreateDeleteTestRequestServer(t, token, "/other")
	defer apiServer.Close()

	err := misc.Delete(apiServer.URL+"/other", authenticator)
	testUtil.IsNil(t, err)
}

func TestPostJSONWithValidResponse(t *testing.T) {
	var apiServer = testUtil.CreatePostJSONTestRequestServer(t, token, `{"id":"id1","name":"Chris"}`, `/endthing`, `{"id":"id1","name":"name"}`)
	defer apiServer.Close()
	actual := TestStruct{}
	ti := TestStruct{ID: "id1", Name: "name"}

	err := misc.PostJSON(apiServer.URL+"/endthing", authenticator, ti, &actual)
	testUtil.IsNil(t, err)
	expected := TestStruct{ID: "id1", Name: "Chris"}

	testUtil.Equals(t, expected, actual)
}

func TestPostJSONWithValidStatusAndEmptyPayloadShouldNotError(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			return
		}))
	defer apiServer.Close()

	actual := TestStruct{}
	ti := TestStruct{ID: "id1", Name: "name"}

	err := misc.PostJSON(apiServer.URL, authenticator, ti, &actual)
	testUtil.IsNil(t, err)

	testUtil.Equals(t, TestStruct{}, actual)
}

func TestGetJSONWithErrorStatusCode(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			return
		}))

	defer apiServer.Close()
	output := TestStruct{}
	err := misc.GetJSON(apiServer.URL, authenticator, &output)
	VerifyNotFoundHTTPStatusInErr(t, err)
}

func TestDeleteWithErrorStatusCode(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			return
		}))

	defer apiServer.Close()
	err := misc.Delete(apiServer.URL, authenticator)
	VerifyNotFoundHTTPStatusInErr(t, err)
}

func TestPutJSONWithErrorStatusCode(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			return
		}))

	defer apiServer.Close()
	output := TestStruct{}
	err := misc.GetJSON(apiServer.URL, authenticator, &output)
	VerifyNotFoundHTTPStatusInErr(t, err)
}

func TestPostJSONWithErrorStatusCode(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			return
		}))

	defer apiServer.Close()
	output := TestStruct{}
	err := misc.PostJSON(apiServer.URL, authenticator, &output, nil)
	VerifyNotFoundHTTPStatusInErr(t, err)
}

func TestGetBody(t *testing.T) {
	reader := testUtil.NewTestReadCloser([]byte(`{"prop" : "val"}`))
	response := http.Response{
		Header:     http.Header{},
		StatusCode: 400,
		Body:       reader,
	}
	httpStatus := misc.NewHTTPStatus(&response, "message")
	payload, err := httpStatus.GetBody()
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "{\"prop\" : \"val\"}", string(payload))

	reader.ErrorOnCallRead = true
	// call the second time and ensure read doesn't occur on the reader
	// and same value is returned.
	payload, err = httpStatus.GetBody()
	testUtil.IsNil(t, err)
	testUtil.Equals(t, "{\"prop\" : \"val\"}", string(payload))

}

func VerifyNotFoundHTTPStatusInErr(t *testing.T, err error) {
	httpStatus, ok := err.(misc.HTTPStatus)
	testUtil.Assert(t, ok, "Expected cast to HTTPStatus")
	testUtil.Equals(t, httpStatus.StatusCode, 404)
}

func TestCallAPI(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.WriteHeader(200) //ok
		}))
	zeroByte := &([]byte{})
	if _, err := misc.CallAPI("HEAD", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := misc.CallAPI("DELETE", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := misc.CallAPI("POST", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}

func TestCallAPIGetContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
			w.Write(body)
		}))
	var resp *http.Response
	if resp, err = misc.CallAPI("GET", apiServer.URL, &fContent, "X-Auth-Token", tokn,
		"Etag", "md5hash-blahblah"); err != nil {
		t.Error(err)
	}
	if strconv.Itoa(len(fContent)) != resp.Header.Get("Content-Length") {
		t.Error(errors.New("Failed: Content-Length comparison"))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(fContent, body) {
		t.Error(errors.New("Failed: Content body comparison"))
	}
}

func TestCallAPIPutContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if strconv.Itoa(len(fContent)) != r.Header.Get("Content-Length") {
				t.Error(errors.New("Failed: Content-Length comparison"))
			}
			if !bytes.Equal(fContent, body) {
				t.Error(errors.New("Failed: Content body comparison"))
			}
			w.WriteHeader(200)
		}))
	if _, err = misc.CallAPI("PUT", apiServer.URL, &fContent, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}

func TestStrcat(t *testing.T) {

	str1 := "abc"
	str2 := "123"
	str3 := ""
	str4 := "`~!@#$%^&*()_+}{[]||';:/?,.<> "
	str5 := "XYZ"
	expectedStr := "abc123`~!@#$%^&*()_+}{[]||';:/?,.<> XYZ"
	result := misc.Strcat(str1, str2, str3, str4, str5)
	if result != expectedStr {
		t.Error(fmt.Errorf("Failed: Strcat doesn't concatenate strings correctly. Expected result = '%s', Actual result = '%s'", expectedStr, result))
	}
}

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
