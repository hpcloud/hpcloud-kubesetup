package machinemetadata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.openstack.org/stackforge/golang-client.git/misc"
)

var openstackMachineMetadataRequestURL = "http://169.254.169.254/openstack/2013-10-17/meta_data.json"

// Response has properties of the machines metadata
type Response struct {
	RandomSeed       string            `json:"random_seed"`
	UUID             string            `json:"uuid"`
	AvailabilityZone string            `json:"availability_zone"`
	HostName         string            `json:"hostname"`
	LaunchIndex      int               `json:"launch_index"`
	PublicKeys       map[string]string `json:"public_keys"`
	Name             string            `json:"name"`
}

// CurrentMachineMetadata returns the machines openstack metadata.
func CurrentMachineMetadata() (Response, error) {
	m, err := getMetadata(openstackMachineMetadataRequestURL)
	return m, err
}

func getMetadata(openStackMetadataURL string) (Response, error) {
	metadata := Response{}
	response, err := http.Get(openStackMetadataURL)
	if err != nil {
		return metadata, err
	}

	if response.StatusCode != http.StatusOK {
		return metadata, misc.NewHTTPStatus(response, fmt.Sprintf("Unable to retrieve metadata. Unexpected status code %d", response.StatusCode))
	}

	err = json.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, err
}
