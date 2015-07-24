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

package requester

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// ExtractSendRequestFunction will get the SendRequestFunction
// if the val implements the requester.Manager interface
func ExtractSendRequestFunction(val interface{}) SendRequestFunction {
	manager, ok := val.(Manager)
	if ok {
		return manager.Function()
	}

	return nil
}

// StandardHTTPRequestMakerGenerator generates a function that will excute
// an http request using the specified Client
func StandardHTTPRequestMakerGenerator(client http.Client) SendRequestFunction {
	return func(request *http.Request) (*http.Response, error) {
		return client.Do(request)
	}
}

// DebugRequestMakerGenerator generates a function that will if debug = true
// writes the http requests and responses to Standard out.
func DebugRequestMakerGenerator(optionalExecutingFunction SendRequestFunction, client *http.Client, debug bool) SendRequestFunction {
	return func(request *http.Request) (response *http.Response, err error) {
		if debug {
			// Note: Use of the dumpRequest after the request has been sent
			// ends up not writing out the payload. Fix is to dump prior to the
			// request being sent.
			dumpRequestToStdOut(request)
		}

		if optionalExecutingFunction == nil {
			if client != nil {
				response, err = client.Do(request)
			} else {
				response, err = http.DefaultClient.Do(request)
			}
		} else {
			response, err = optionalExecutingFunction(request)
		}

		if debug {
			dumpResponseToStdOut(response)
		}

		return response, err
	}
}

func dumpResponseToStdOut(response *http.Response) {

	defer func() {
		if error := recover(); error != nil {
			fmt.Fprintln(os.Stdout, "Error: panic occurred creating an http response dump")
		}
	}()

	dumpResponseBytes, err := httputil.DumpResponse(response, true)
	if err == nil {
		fmt.Fprintln(os.Stdout, string(dumpResponseBytes))
		fmt.Println() // emptyline to separate from response better.
	}
}

func dumpRequestToStdOut(request *http.Request) {
	dumpRequestBytes, err := httputil.DumpRequest(request, true)
	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------------")
	if err == nil {
		fmt.Fprintln(os.Stdout, string(dumpRequestBytes))
	}
}
