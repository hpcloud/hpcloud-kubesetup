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

package misc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
)

// The default HTTP Transport. Users are allowed to override
// the transport so they can specify parameters for their own
// use-cases.
var transport = &http.Transport{}

// HTTPStatus is the error that that is returned when an
// error status code occurs. Allows developers to use this information
// in understanding the type of error that occurs for a particular request.
type HTTPStatus struct {
	StatusCode         int
	Message            string
	Header             http.Header
	responseBodyReader io.Reader
	savedBody          []byte
	savedBodyError     error
	bodyRead           bool
}

// NewHTTPStatus with specified values.
func NewHTTPStatus(response *http.Response, errorMessage string) HTTPStatus {
	return HTTPStatus{StatusCode: response.StatusCode, Message: errorMessage, Header: response.Header, responseBodyReader: response.Body}
}

// GetBody gets the full body of the response.
func (httpStatus *HTTPStatus) GetBody() ([]byte, error) {
	if !httpStatus.bodyRead {
		// If the status code indicates success then just return an internal assertion error as this case shouldn't be hit
		if httpStatus.StatusCode >= 200 && httpStatus.StatusCode <= 300 {
			return []byte{}, fmt.Errorf("Error, internal assert, Response has already been read for a non-error status code")
		}

		if httpStatus.responseBodyReader != nil {
			httpStatus.savedBody, httpStatus.savedBodyError = ioutil.ReadAll(httpStatus.responseBodyReader)
		}

		httpStatus.bodyRead = true
	}

	return httpStatus.savedBody, httpStatus.savedBodyError
}

// Error ensures that HTTPStatus implements the Error interface.
func (httpStatus HTTPStatus) Error() string {
	return httpStatus.Message
}

// Transport allows a custom net/http.Transport to be specified
func Transport(t *http.Transport) error {
	transport = t

	return nil
}

// NewHTTPClient creates an HTTP client using our transport for SSL purposes.
func NewHTTPClient() *http.Client {
	return &http.Client{Transport: transport}
}

// CallGetAPI invokes HTTP GET request.
func CallGetAPI(url string, h ...string) (header http.Header, responseByteArr []byte, err error) {
	resp, err := CallAPI("GET", url, zeroByte, h...)
	header = resp.Header

	if _, err = CheckHTTPResponseStatusCode(resp); err != nil {
		return header, nil, err
	}

	responseByteArr, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return header, responseByteArr, err
	}

	return header, responseByteArr, nil
}

//CallAPI sends an HTTP request using "method" to "url".
//For uploading / sending file, caller needs to set the "content".  Otherwise,
//set it to zero length []byte. If Header fields need to be set, then set it in
// "h".  "h" needs to be even numbered, i.e. pairs of field name and the field
//content.
//
//fileContent, err := ioutil.ReadFile("fileName.ext");
//
//resp, err := CallAPI("PUT", "http://domain/hello/", &fileContent,
//"Name", "world")
//
//is similar to: curl -X PUT -H "Name: world" -T fileName.ext
//http://domain/hello/
func CallAPI(method, url string, content *[]byte, h ...string) (*http.Response, error) {
	if len(h)%2 == 1 { //odd #
		return nil, errors.New("syntax err: # header != # of values")
	}
	//I think the above err check is unnecessary and wastes cpu cycle, since
	//len(h) is not determined at run time. If the coder puts in odd # of args,
	//the integration testing should catch it.
	//But hey, things happen, so I decided to add it anyway, although you can
	//comment it out, if you are confident in your test suites.
	var req *http.Request
	var err error
	contentLength := int64(len(*content))
	if contentLength > 0 {
		req, err = http.NewRequest(method, url, bytes.NewReader(*content))
		//req.Body = *(new(io.ReadCloser)) //these 3 lines do not work but I am
		//req.Body.Read(content)           //keeping them here in case I wonder why
		//req.Body.Close()                 //I did not implement it this way :)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	req.ContentLength = contentLength

	if err != nil {
		return nil, err
	}
	for i := 0; i < len(h)-1; i = i + 2 {
		req.Header.Set(h[i], h[i+1])
	}

	resp, err := makeRequest(nil, req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// Delete sends an Http Request with using the "DELETE" method and with
// an "X-Auth-Token" header set to the specified token value. The request
// is made by the specified client.
func Delete(url string, authenticator common.Authenticator) (err error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return err
	}

	return DeleteWithTokenAndRequester(url, token, r)
}

// DeleteWithTokenAndRequester sends an Http Request with using the "DELETE" method and with
// an "X-Auth-Token" header set to the specified token value. The request
// is made by the specified client.
func DeleteWithTokenAndRequester(url string, token string, r requester.SendRequestFunction) (err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token)

	resp, err := makeRequest(r, req)
	if err != nil {
		return err
	}

	// Expecting a successful delete
	if !(resp.StatusCode == 200 || resp.StatusCode == 202 || resp.StatusCode == 204) {
		err = NewHTTPStatus(resp, fmt.Sprintf("Unexpected server response status code on Delete '%d'", resp.StatusCode))
		return
	}

	return nil
}

// Get sends an HTTP Request with using the "GET" method and with
// an "Accept" header set to "*/*" and the authentication token
// set to the specified token value. The request is made by the
// specified client. The response is a straight net/http.Response
// The caller is responsible for closing the response body!
func Get(url string, authenticator common.Authenticator) (*http.Response, error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Auth-Token", token)

	resp, err := makeRequest(r, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 200 && resp.StatusCode != 300 {
		err = HTTPStatus{StatusCode: resp.StatusCode, Message: "Error: status code != 200, 201, 202, or 300, actual status code '" + resp.Status + "'"}
	}

	return resp, err
}

// Head sends an HTTP Request with using the "HEAD" method and with
// an "Accept" header set to "*/*" and the authentication token
// set to the specified token value. The request is made by the
// specified client. The response is a straight net/http.Response
// The caller is responsible for closing the response body!
func Head(url string, authenticator common.Authenticator) (*http.Response, error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Auth-Token", token)

	resp, err := makeRequest(r, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 200 && resp.StatusCode != 300 {
		err = HTTPStatus{StatusCode: resp.StatusCode, Message: "Error: status code != 200, 201, 202, or 300, actual status code '" + resp.Status + "'"}
	}

	return resp, err
}

//GetJSON sends an Http Request with using the "GET" method and with
//an "Accept" header set to "application/json" and the authentication token
//set to the specified token value. The request is made by the
//specified client. The val interface should be a pointer to the
//structure that the json response should be decoded into.
func GetJSON(url string, authenticator common.Authenticator, val interface{}) (err error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return err
	}

	return GetJSONWithTokenAndRequester(url, token, r, val)
}

//GetJSONWithTokenAndRequester sends an Http Request with using the "GET" method and with
//an "Accept" header set to "application/json" and the authentication token
//set to the specified token value. The request is made by the
//specified client. The val interface should be a pointer to the
//structure that the json response should be decoded into.
func GetJSONWithTokenAndRequester(url string, token string, requester requester.SendRequestFunction, val interface{}) (err error) {
	req, err := createJSONGetRequest(url, token)
	if err != nil {
		return err
	}

	err = executeRequestCheckStatusDecodeJSONResponse(requester, req, val)
	if err != nil {
		return err
	}

	return nil
}

// PostJSON sends an Http Request with using the "POST" method and with
// a "Content-Type" header with application/json and "X-Auth-Token" header
// if a value is specified. The inputValue is encoded to json
// and sent in the body of the request. The response json body is
// decoded into the outputValue. If the response does sends an invalid
// or error status code then an error will be returned.
func PostJSON(url string, authenticator common.Authenticator, inputValue interface{}, outputValue interface{}) (err error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return err
	}

	return PostJSONWithTokenAndRequester(url, token, r, inputValue, outputValue)
}

// PostJSONWithTokenAndRequester sends an Http Request with using the "POST" method and with
// a "Content-Type" header with application/json and "X-Auth-Token" header
// if a value is specified. The inputValue is encoded to json
// and sent in the body of the request. The response json body is
// decoded into the outputValue. If the response does sends an invalid
// or error status code then an error will be returned.
func PostJSONWithTokenAndRequester(url string, token string, requestor requester.SendRequestFunction, inputValue interface{}, outputValue interface{}) (err error) {
	req, err := createPostRequest(url, token, inputValue)
	if err != nil {
		return err
	}

	resp, err := makeRequest(requestor, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 200 && resp.StatusCode != 300 {
		err = NewHTTPStatus(resp, "Error: status code != 200, 201, 202, or 300, actual status code '"+resp.Status+"'")
		return
	}

	// Empty content should just return without decoding as there is no content.
	if resp.ContentLength == 0 {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&outputValue)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func createPostRequest(url string, token string, inputValue interface{}) (*http.Request, error) {
	var req *http.Request
	var err error
	if inputValue != nil {
		body, err := json.Marshal(inputValue)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "gzip,deflate")
	// TODO: Need to add the specified char set as utf8 as this is what is produced by golang.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}

	return req, nil
}

// PutJSON sends an Http Request with using the "PUT" method and with
// a "Content-Type" header with application/json and X-Auth-Token" header
// set to the specified token value. The inputValue is encoded to json
// and sent in the body of the request. The response json body is
// decoded into the outputValue. If the response does sends a
// non 200 status code then an error will be returned.
func PutJSON(url string, authenticator common.Authenticator, inputValue interface{}, outputValue interface{}) (err error) {
	r := requester.ExtractSendRequestFunction(authenticator)
	token, err := authenticator.GetToken()
	if err != nil {
		return err
	}

	return PutJSONWithTokenAndRequester(url, token, r, inputValue, outputValue)
}

// PutJSONWithTokenAndRequester sends an Http Request with using the "PUT" method and with
// a "Content-Type" header with application/json and X-Auth-Token" header
// set to the specified token value. The inputValue is encoded to json
// and sent in the body of the request. The response json body is
// decoded into the outputValue. If the response does sends a
// non 200 status code then an error will be returned.
func PutJSONWithTokenAndRequester(url string, token string, requestor requester.SendRequestFunction, inputValue interface{}, outputValue interface{}) (err error) {
	body, err := json.Marshal(inputValue)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, err := makeRequest(requestor, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = NewHTTPStatus(resp, "Error: status code != 200 actual status code '"+resp.Status+"'")
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&outputValue)
	defer resp.Body.Close()

	return err
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204.
// Otherwise return error message.
func CheckHTTPResponseStatusCode(resp *http.Response) (HTTPStatus, error) {
	var message string
	switch resp.StatusCode {
	case 200, 201, 202, 204, 300:
		return HTTPStatus{StatusCode: resp.StatusCode}, nil
	case 400:
		message = "Error: response == 400 bad request"
	case 401:
		message = "Error: response == 401 unauthorised"
	case 403:
		message = "Error: response == 403 forbidden"
	case 404:
		message = "Error: response == 404 not found"
	case 405:
		message = "Error: response == 405 method not allowed"
	case 409:
		message = "Error: response == 409 conflict"
	case 413:
		message = "Error: response == 413 over limit"
	case 415:
		message = "Error: response == 415 bad media type"
	case 422:
		message = "Error: response == 422 unprocessable"
	case 429:
		message = "Error: response == 429 too many request"
	case 500:
		message = "Error: response == 500 instance fault / server err"
	case 501:
		message = "Error: response == 501 not implemented"
	case 503:
		message = "Error: response == 503 service unavailable"
	default:
		message = "Error: unexpected response status code"
	}

	httpStatus := NewHTTPStatus(resp, message)

	return httpStatus, httpStatus
}

// ContentTypeIsJSON determines if the content-type is application/json
func ContentTypeIsJSON(header http.Header) bool {
	contentType := header.Get("Content-Type")
	if contentType != "application/json" {
		return false
	}

	return true
}

// Strcat concatenates all input strings.
// Refer the following link for "Fastest string contatenation":
// http://golang-examples.tumblr.com/post/86169510884/fastest-string-contatenation
func Strcat(strs ...string) string {
	var buffer bytes.Buffer
	for _, s := range strs {
		buffer.WriteString(s)
	}
	return buffer.String()
}

///////////////////////////////////////////////////////////////////////////////
// Private Structs
///////////////////////////////////////////////////////////////////////////////

var zeroByte = new([]byte) //pointer to empty []byte

type readCloser struct {
	io.Reader
}

///////////////////////////////////////////////////////////////////////////////
// Private Functions
///////////////////////////////////////////////////////////////////////////////

func (readCloser) Close() error {
	//cannot put this func inside CallAPI; golang disallow nested func
	return nil
}

func makeRequest(r requester.SendRequestFunction, request *http.Request) (*http.Response, error) {
	if r == nil {
		client := NewHTTPClient()
		r = requester.StandardHTTPRequestMakerGenerator(*client)
	}

	return r(request)
}

func createJSONGetRequest(url string, token string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", token)

	return req, nil
}

func executeRequestCheckStatusDecodeJSONResponse(r requester.SendRequestFunction, req *http.Request, val interface{}) (err error) {
	resp, err := makeRequest(r, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 200 && resp.StatusCode != 300 {
		err = HTTPStatus{StatusCode: resp.StatusCode, Message: "Error: status code != 200, 201, 202, or 300, actual status code '" + resp.Status + "'"}
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&val)
	defer resp.Body.Close()

	return err
}
