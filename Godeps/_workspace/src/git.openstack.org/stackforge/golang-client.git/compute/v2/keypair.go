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

// KeyPairResponse is a structure for all properties of for a tenant
type KeyPairResponse struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	FingerPrint string `json:"fingerprint"`
	UserID      string `json:"user_id"`
}

// KeyPairs will issue a get request to retrieve the all keypairs.
func (computeService Service) KeyPairs() ([]KeyPairResponse, error) {
	var kp = keyPairsResponseContainer{}
	url, err := computeService.buildRequestURL("/os-keypairs")
	if err != nil {
		return nil, err
	}

	err = misc.GetJSON(url, computeService.authenticator, &kp)
	if err != nil {
		return nil, err
	}

	var keypairs []KeyPairResponse
	for _, keyPairContainedItem := range kp.KeyPairs {
		keypairs = append(keypairs, keyPairContainedItem.KeyPair)
	}

	return keypairs, nil
}

// KeyPair will issue a get request to retrieve the specified keypair.
func (computeService Service) KeyPair(name string) (KeyPairResponse, error) {
	var kp = keyPairResponseContainer{}
	url, err := computeService.buildRequestURL("/os-keypairs/", name)
	if err != nil {
		return kp.KeyPair, err
	}

	err = misc.GetJSON(url, computeService.authenticator, &kp)
	return kp.KeyPair, err
}

// CreateKeyPair will send a POST request to create a new keypair with the specified parameters.
func (computeService Service) CreateKeyPair(name string, publickey string) (KeyPairResponse, error) {
	createKeypairContainerValues := keyPairRequestContainer{keyPairRequest{Name: name, PublicKey: publickey}}
	outKeyContainer := keyPairResponseContainer{}
	reqURL, err := computeService.buildRequestURL("/os-keypairs")
	if err != nil {
		return outKeyContainer.KeyPair, err
	}

	err = misc.PostJSON(reqURL, computeService.authenticator, createKeypairContainerValues, &outKeyContainer)
	return outKeyContainer.KeyPair, err
}

// DeleteKeyPair will delete the keypair.
func (computeService Service) DeleteKeyPair(name string) (err error) {
	reqURL, err := computeService.buildRequestURL("/os-keypairs/", name)
	if err != nil {
		return err
	}

	return misc.Delete(reqURL, computeService.authenticator)
}

type keyPairsResponseContainer struct {
	KeyPairs []keyPairResponseContainer `json:"keypairs"`
}

type keyPairResponseContainer struct {
	KeyPair KeyPairResponse `json:"keypair"`
}

type keyPairRequestContainer struct {
	KeyPair keyPairRequest `json:"keypair"`
}

type keyPairRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}
