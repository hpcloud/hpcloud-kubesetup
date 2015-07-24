package machinemetadata

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/misc"
	"github.com/stretchr/testify/assert"
)

var sampleValue = `{
  "random_seed": "b925GxNd",
  "uuid": "4d6e242a-ba48-4a95-b062-0e669839035a",
  "availability_zone": "az2",
  "hostname": "workvm.novalocal",
  "launch_index": 0,
  "public_keys": {
    "Macbook": "ssh-rsa BBBB3NzaC1yc2EAAASnxHC4Pw"
  },
  "name": "workVM"
}`

func TestGetOpenStackMetadata(t *testing.T) {
	assert := assert.New(t)
	testSvr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Here")
			assert.Equal("GET", r.Method)
			w.Write([]byte(sampleValue))
			w.WriteHeader(http.StatusOK)
		}))
	defer testSvr.Close()
	metadata, err := getMetadata(testSvr.URL)
	assert.Nil(err)

	assert.Equal("az2", metadata.AvailabilityZone)
	assert.Equal("workvm.novalocal", metadata.HostName)
	assert.Equal(0, metadata.LaunchIndex)
	assert.Equal("workVM", metadata.Name)
	assert.Equal(map[string]string{"Macbook": "ssh-rsa BBBB3NzaC1yc2EAAASnxHC4Pw"}, metadata.PublicKeys)
	assert.Equal("b925GxNd", metadata.RandomSeed)
	assert.Equal("4d6e242a-ba48-4a95-b062-0e669839035a", metadata.UUID)
}

func TestGetOpenStackMetadataErrorOnRequest(t *testing.T) {
	assert := assert.New(t)
	testSvr := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal("GET", r.Method)
			w.WriteHeader(http.StatusBadRequest)
		}))
	defer testSvr.Close()
	_, err := getMetadata(testSvr.URL)
	assert.NotNil(err)
	httpStatus, isType := err.(misc.HTTPStatus)
	assert.Equal(true, isType)
	assert.Equal(httpStatus.StatusCode, 400)
	assert.Equal("Unable to retrieve metadata. Unexpected status code 400", err.Error())
}
