package identity_test

import (
	"fmt"
	"os"
	"testing"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	database "git.openstack.org/stackforge/golang-client.git/database/v1"
	"git.openstack.org/stackforge/golang-client.git/identity/common"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	image "git.openstack.org/stackforge/golang-client.git/image/v1"
	"git.openstack.org/stackforge/golang-client.git/misc/requester"
	network "git.openstack.org/stackforge/golang-client.git/network/v2"
	objectstorage "git.openstack.org/stackforge/golang-client.git/objectstorage/v1"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var initializationMap = map[string]map[string]string{
	"PublicHelion": {
		"compute":      "2",
		"network":      "2",
		"image":        "1",
		"object-store": "1",
		"database":     "v13.6"},
	"PrivateHelion": {
		"compute":      "2",
		"network":      "2",
		"image":        "1",
		"object-store": "1",
		/*"database": "1" Not running because endpoint doesn't work*/},
}

var currentRunningMap = map[string]string{}

func init() {
	openStackType := os.Getenv("OpenStackType")
	if openStackType == "" {
		// Initialize to public helion by default.
		currentRunningMap = initializationMap["PublicHelion"]
	} else {
		currentRunningMap = initializationMap[openStackType]
	}
}

func TestComputeV2ServiceCanQuery(t *testing.T) {
	_, _, _, authenticator := getServiceURLInitialize(t, "compute")
	computeService := compute.NewService(authenticator)
	_, err := computeService.Flavors()
	testUtil.IsNil(t, err)
}

func TestNetworkV2ServiceCanQuery(t *testing.T) {
	_, _, _, authenticator := getServiceURLInitialize(t, "network")
	networkService := network.NewService(authenticator)
	_, err := networkService.Routers()
	testUtil.IsNil(t, err)
}

func TestImageV1ServiceCanQuery(t *testing.T) {
	_, _, _, authenticator := getServiceURLInitialize(t, "image")
	imageService := image.NewService(authenticator)
	_, err := imageService.Images()
	testUtil.IsNil(t, err)
}

func TestObjectStorageV1ServiceCanQuery(t *testing.T) {
	tokenID, serviceURL, _, _ := getServiceURLInitialize(t, "object-store")

	_, err := objectstorage.ListContainers(0, "", serviceURL+"?format=json", tokenID)
	testUtil.IsNil(t, err)
}

func TestDatabaseV1ServiceCanQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	version := skipIfServiceNotFound(t, "database")

	_, authenticator, ap := authenticate(t)
	// Either it works or fails with the failure below.
	_, err := authenticator.GetServiceURL("database", version)
	if err != nil {
		if err.Error() == "ServiceCatalog does not contain serviceType 'database'" {
			t.Skip("Database service does not exist in region, skipping")
		} else {
			testUtil.Equals(t, fmt.Sprintf(identity.ServiceEndpointWithSpecifiedRegionAndVersionNotFound, "database", ap.Region, "v13.6"), err.Error())
		}

		return
	}

	databaseService := database.NewService(authenticator)
	_, err = databaseService.Flavors()
	testUtil.IsNil(t, err)
}

func getServiceURLInitialize(t *testing.T, serviceType string) (string, string, common.AuthenticationParameters, common.Authenticator) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	version, ok := currentRunningMap[serviceType]
	if !ok {
		t.Skip("skipping test as the service type is not in the map specifying the version to use.")
	}

	tokenID, authenticator, ap := authenticate(t)
	serviceURL, err := authenticator.GetServiceURL(serviceType, version)
	testUtil.IsNil(t, err)
	testUtil.Assert(t, serviceURL != "", "No service URL was found.")

	return tokenID, serviceURL, ap, authenticator
}

func skipIfServiceNotFound(t *testing.T, serviceType string) string {
	version, ok := currentRunningMap[serviceType]
	if !ok {
		t.Skip("skipping test as the service type '%s' is not in the map specifying the version to use.", serviceType)
	}

	return version
}

func authenticate(t *testing.T) (string, common.Authenticator, common.AuthenticationParameters) {
	ap := getAuthenticationParameters(t)
	authenticator := identity.Authenticate(ap)
	authenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))

	tokenID, err := authenticator.GetToken()
	testUtil.IsNil(t, err)
	testUtil.Assert(t, tokenID != "", "No tokenID was found.")

	return tokenID, &authenticator, ap
}

func getAuthenticationParameters(t *testing.T) common.AuthenticationParameters {
	ap, err := common.FromEnvVars()
	if err != nil {
		t.Fatal("Integration test cannot proceed because env vars are not set.")
	}

	return ap
}

func TestAuthenticate(t *testing.T) {
	//Authenticate with Global envs
	ap := getAuthenticationParameters(t)

	authenticator := identity.Authenticate(ap)
	authenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))

	tokenID, err := authenticator.GetToken()
	if tokenID == "" || err != nil {
		t.Fatal("Error while authenticating")
	}

	tokenBasedAp := common.AuthenticationParameters{
		AuthURL:    ap.AuthURL,
		TenantID:   ap.TenantID,
		TenantName: ap.TenantName,
		Region:     ap.Region,
		CACert:     ap.CACert,
		AuthToken:  tokenID,
	}

	tokenBasedAuthenticator := identity.Authenticate(tokenBasedAp)
	tokenBasedAuthenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, nil, testing.Verbose()))

	otherToken, err := tokenBasedAuthenticator.GetToken()
	if otherToken == "" || err != nil {
		t.Fatal("Error while authenticating with token")
	}

	//Test if able to connect to compute, and get Flavors
	serviceURL, err := tokenBasedAuthenticator.GetServiceURL("compute", "2")
	if serviceURL == "" || err != nil {
		t.Fatal("Error while fetching Compute Service with Token")
	}

	computeService := compute.NewService(&tokenBasedAuthenticator)
	_, err = computeService.Flavors()
	if err != nil {
		t.Fatal("Error while fetching Compute flavors with Token")
	}
}
