package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	compute "git.openstack.org/stackforge/golang-client.git/compute/v2"
	common "git.openstack.org/stackforge/golang-client.git/identity/common"
	identity "git.openstack.org/stackforge/golang-client.git/identity/v2"
	image "git.openstack.org/stackforge/golang-client.git/image/v1"
	misc "git.openstack.org/stackforge/golang-client.git/misc"
	requester "git.openstack.org/stackforge/golang-client.git/misc/requester"
	network "git.openstack.org/stackforge/golang-client.git/network/v2"

	"github.com/codegangsta/cli"
	"github.com/parnurzeal/gorequest"

	"gopkg.in/yaml.v2"
)

var version = "0.0.3"

var (
	computeService compute.Service
	networkService network.Service
	imageService   image.Service
	config         configContainer
	keypair        compute.KeyPairResponse
	netwrk         network.Response
	subnets        []network.SubnetResponse
	servers        []compute.Server
	ports          []network.PortResponse
	flavorMap      map[string]string
)

type configContainer struct {
	Nodes            map[string]configNode `yaml:"hosts"`
	SSHKey           string                `yaml:"sshkey"`
	Network          string                `yaml:"network"`
	AvailabilityZone string                `yaml:"availabilityZone"`
	OrderedNodeKeys  []string
}

type configNode struct {
	IP       string `yaml:"ip"`
	IsMaster bool   `yaml:"ismaster"`
	VMImage  string `yaml:"vm-image"`
	VMSize   string `yaml:"vm-size"`
	ServerID string
}

func main() {

	app := cli.NewApp()
	app.Name = "hpcloud-kubesetup"
	app.Usage = "Kubernetes cluster setup for HP Helion OpenStack"
	app.Version = version
	app.Author = "Gert Drapers"
	app.Email = "gert.drapers@hp.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  Config,
			Value: DefaultConfig,
			Usage: "Kubernetes cluster configuration file",
		},
		cli.StringFlag{
			Name:   AuthURL,
			Value:  "",
			Usage:  "OpenStack authentication URL",
			EnvVar: AuthURLEnv,
		},
		cli.StringFlag{
			Name:   TenantID,
			Value:  "",
			Usage:  "OpenStack tenant id",
			EnvVar: TenantIDEnv,
		},
		cli.StringFlag{
			Name:   TenantName,
			Value:  "",
			Usage:  "OpenStack tenant name",
			EnvVar: TenantNameEnv,
		},
		cli.StringFlag{
			Name:   Username,
			Value:  "",
			Usage:  "OpenStack user name",
			EnvVar: UsernameEnv,
		},
		cli.StringFlag{
			Name:   Password,
			Value:  "",
			Usage:  "OpenStack password",
			EnvVar: PasswordEnv,
		},
		cli.StringFlag{
			Name:   RegionName,
			Value:  "",
			Usage:  "OpenStack region name",
			EnvVar: RegionNameEnv,
		},
		cli.StringFlag{
			Name:   AuthToken,
			Value:  "",
			Usage:  "OpenStack auth token",
			EnvVar: AuthTokenEnv,
		},
		cli.StringFlag{
			Name:   CACert,
			Value:  "",
			Usage:  "OpenStack TLS (https) server certificate",
			EnvVar: CACertEnv,
		},
		cli.BoolFlag{
			Name:  SkipSSLValidation,
			Usage: "Skip SSL validation",
		},
		cli.BoolFlag{
			Name:  Debug,
			Usage: "Enable debug spew",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   Install,
			Usage:  "Create Kubernetes cluster",
			Action: installAction,
		},
		/*
			{
				Name:   Status,
				Usage:  "Status of Kubernetes cluster",
				Action: statusAction,
			},
		*/
		{
			Name:   Uninstall,
			Usage:  "Remove Kubernetes cluster",
			Action: uninstallAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(1)
}

func installAction(c *cli.Context) {

	initTask(c)
	uninstallTask(c)
	createCloudConfigTask(c)
	installTask(c)
	statusTask(c)
	assignIPAddressTask(c)
}

func statusAction(c *cli.Context) {

	initTask(c)
	statusTask(c)
}

func uninstallAction(c *cli.Context) {

	initTask(c)
	uninstallTask(c)
}

func initTask(c *cli.Context) {

	var err error

	config, err = readConfigFile(c.GlobalString(Config))
	if err != nil {
		log.Fatal(err.Error())
	}
	for k := range config.Nodes {
		config.OrderedNodeKeys = append(config.OrderedNodeKeys, k)
	}
	sort.Strings(config.OrderedNodeKeys)

	config.Log()

	log.Printf("%-20s - %s\n", AuthURLEnv, c.GlobalString(AuthURL))
	log.Printf("%-20s - %s\n", TenantIDEnv, c.GlobalString(TenantID))
	log.Printf("%-20s - %s\n", TenantNameEnv, c.GlobalString(TenantName))
	log.Printf("%-20s - %s\n", UsernameEnv, c.GlobalString(Username))
	log.Printf("%-20s - %s\n", RegionNameEnv, c.GlobalString(RegionName))
	log.Printf("%-20s - %s\n", AuthTokenEnv, c.GlobalString(AuthToken))
	log.Printf("%-20s - %s\n", CACertEnv, c.GlobalString(CACert))
	log.Printf("%-20s - %v\n", SkipSSLValidation, c.GlobalBool(SkipSSLValidation))
	log.Printf("%-20s - %v\n", Debug, c.GlobalBool(Debug))

	authParameters := common.AuthenticationParameters{
		AuthURL:    c.GlobalString(AuthURL),
		Username:   c.GlobalString(Username),
		Password:   c.GlobalString(Password),
		Region:     c.GlobalString(RegionName),
		TenantID:   c.GlobalString(TenantID),
		TenantName: c.GlobalString(TenantName),
		CACert:     c.GlobalString(CACert),
	}

	var transport *http.Transport

	if c.GlobalString(CACert) != "" {

		pemData, err := ioutil.ReadFile(c.GlobalString(CACert))
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "Unable to load CA Certificate file", c.GlobalString(CACert)))
		}

		certPool := x509.NewCertPool()

		if !certPool.AppendCertsFromPEM(pemData) {
			log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "Invalid CA Certificates in file", c.GlobalString(CACert)))
		}

		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: c.GlobalBool(SkipSSLValidation),
			},
		}
		misc.Transport(transport)

	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.GlobalBool(SkipSSLValidation),
			},
		}
		misc.Transport(transport)
	}

	authenticator := identity.Authenticate(authParameters)
	authenticator.SetFunction(requester.DebugRequestMakerGenerator(nil, &http.Client{Transport: transport}, c.GlobalBool(Debug)))

	token, err := authenticator.GetToken()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}
	log.Printf("%-20s - %s\n", "token", token)

	computeService = compute.NewService(&authenticator)

	networkService = network.NewService(&authenticator)

	imageService = image.NewService(&authenticator)

	keypair, err = computeService.KeyPair(config.SSHKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s %s\n", "error:", "get keypair", config.SSHKey, err.Error()))
	}

	var q = network.QueryParameters{Name: config.Network}
	networks, err := networkService.QueryNetworks(q)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s %s\n", "error:", "get network by name", config.Network, err.Error()))
	}

	netwrk, err = networkService.Network(networks[0].ID)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s %s\n", "error:", "getting network by id", networks[0].ID, err.Error()))
	}
	log.Printf("%-20s - %s\n", "network", netwrk.ID)

	subnets, err = networkService.Subnets()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "get subnets", err.Error()))
	}

	ports, err = networkService.Ports()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "get ports", err.Error()))
	}

	sort.Sort(PortByName(ports))

	servers, err = computeService.Servers()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "get servers", err.Error()))
	}

	sort.Sort(ServerByName(servers))

	flavors, err := computeService.Flavors()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "get flavors", err.Error()))
	}

	flavorMap = make(map[string]string)
	for _, p := range flavors {
		flavorMap[p.Name] = p.ID
	}
}

func uninstallTask(c *cli.Context) {

	for _, v := range servers {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete server", v.Name)

			err := computeService.DeleteServer(v.ID)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete server", v.Name, "COMPLETED")
		}
	}

	for _, v := range ports {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete port", v.Name)

			err := networkService.DeletePort(v.ID)
			if noErrorOn404(err) != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete port", v.Name, "COMPLETED")
		}
	}
}

func createCloudConfigTask(c *cli.Context) {

	nodeID := 0

	masterIP, err := getMasterIP(config.Nodes)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s %s\n", "error:", "get master IP", err.Error()))
	}

	discovery := getDiscoveryKey()

	for k, v := range config.Nodes {

		log.Printf("%-20s - %s\n", "create cloudconfig", k)

		data := make(map[string]string)

		data["filename"] = k + ".yml"

		if v.IsMaster {
			data["role"] = "master"
		} else {
			data["role"] = "node"
		}

		data["discovery"] = discovery
		data["master"] = masterIP
		data["hostname"] = k
		data["ip"] = v.IP
		data["sshkey"] = keypair.PublicKey

		nodeID = nodeID + 1

		createCloudConfig(data)

		log.Printf("%-20s - %s %s\n", "create cloudconfig", data["filename"], "COMPLETED")
	}
}

func installTask(c *cli.Context) {

	for _, v := range config.OrderedNodeKeys {

		log.Printf("%-20s - %s %s\n", "create port", v, config.Nodes[v].IP)

		newPort := network.CreatePortParameters{}
		newPort.Name = v
		newPort.AdminStateUp = true
		newPort.NetworkID = netwrk.ID
		newPort.FixedIPs = []network.FixedIP{{IPAddress: config.Nodes[v].IP, SubnetID: netwrk.Subnets[0]}}

		port, err := networkService.CreatePort(newPort)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		log.Printf("%-20s - %s %s\n", "create port", port.ID, "COMPLETED")

		log.Printf("%-20s - %s %s\n", "create server", v, config.Nodes[v].IP)

		userdata, err := getUserData(v + ".yml")
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		imageQuery := image.QueryParameters{Name: config.Nodes[v].VMImage}
		images, err := imageService.QueryImages(imageQuery)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}
		log.Printf("%-20s - %s\n", "image", images[0].ID)

		log.Printf("%-20s - %s\n", "flavor", flavorMap[config.Nodes[v].VMSize])

		newServer := compute.ServerCreationParameters{}
		newServer.Name = v
		newServer.ImageRef = images[0].ID
		newServer.FlavorRef = flavorMap[config.Nodes[v].VMSize]
		newServer.KeyPairName = keypair.Name
		newServer.UserData = &userdata
		newServer.Networks = []compute.ServerNetworkParameters{{UUID: port.NetworkID, Port: port.ID}}
		newServer.SecurityGroups = []compute.SecurityGroup{{Name: "default"}}
		newServer.AvailabilityZone = &config.AvailabilityZone

		server, err := computeService.CreateServer(newServer)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		log.Printf("%-20s - %s %s\n", "create server", "password", server.AdminPass)
		log.Printf("%-20s - %s %s\n", "create server", server.ID, "COMPLETED")

		node := config.Nodes[v]
		node.ServerID = server.ID
		config.Nodes[v] = node
	}
}

func statusTask(c *cli.Context) {

	for _, v := range config.OrderedNodeKeys {

		prevStatus := ""
		for {
			server, err := computeService.ServerDetail(config.Nodes[v].ServerID)

			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			if prevStatus != server.Status {
				log.Printf("%-20s - %s %s\n", "server status", server.Name, server.Status)
				prevStatus = server.Status
			}

			if server.Status == "ACTIVE" {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func assignIPAddressTask(c *cli.Context) {

	var unAssigned []compute.FloatingIP
	floatingIPs, err := computeService.FloatingIPs()
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	for _, v := range floatingIPs {
		if len(v.InstanceID) == 0 {
			unAssigned = append(unAssigned, v)
		}
	}

	for k, v := range config.Nodes {

		if !v.IsMaster {
			continue
		}

		if len(unAssigned) == 0 {

			log.Printf("%-20s - %s %s\n", "create public IP", "", "")

			var fp compute.FloatingIP
			fp, err := computeService.CreateFloatingIPInDefaultPool()
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}
			unAssigned = append(unAssigned, fp)

			log.Printf("%-20s - %s %s\n", "create public IP", fp.IP, "COMPLETED")

		}

		log.Printf("%-20s - %s %s\n", "associate IP", k, unAssigned[0].IP)

		err := computeService.ServerAction(v.ServerID, "addFloatingIp", "address", unAssigned[0].IP)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		log.Printf("%-20s - %s %s\n", "associate IP", k, "COMPLETED")
		unAssigned = append(unAssigned[1:])
	}
}

func readConfigFile(filename string) (config configContainer, err error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return
	}

	err = nil
	return
}

func (config configContainer) Log() {
	for k, v := range config.Nodes {
		log.Printf("%-20s - %s %v\n", "config file", k, v)
	}
	log.Printf("%-20s - %s %s\n", "config file", "SSHKey", config.SSHKey)
	log.Printf("%-20s - %s %s\n", "config file", "Network", config.Network)
	log.Printf("%-20s - %s %s\n", "config file", "AvailabilityZone", config.AvailabilityZone)

}

func createCloudConfig(data map[string]string) error {

	var b bytes.Buffer
	f, err := os.Create(data["filename"])
	if err != nil {
		return err
	}

	w := io.MultiWriter(f, &b)
	if data["role"] == "master" {
		if err := masterTmpl.Execute(w, data); err != nil {
			fmt.Println("err masterTmpl")
			return err
		}
	} else {
		if err := nodeTmpl.Execute(w, data); err != nil {
			fmt.Println("err nodeTmpl")
			return err
		}
	}

	f.Close()
	return nil
}

func getUserData(filename string) (encodedStr string, err error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	encodedStr = base64.StdEncoding.EncodeToString(b)
	err = nil
	return
}

func getMasterIP(nodeList map[string]configNode) (string, error) {

	for _, v := range nodeList {
		if v.IsMaster {
			return v.IP, nil
		}
	}
	return "", fmt.Errorf("No master IP address found")
}

/*
	CoreOS Cluster Discovery ID
	See https://coreos.com/docs/cluster-management/setup/cluster-discovery/
	for details
*/
func getDiscoveryKey() string {

	req := gorequest.New()

	_, body, errs := req.Get("https://discovery.etcd.io/new").
		Set("Content-Type", "text/plain").
		Set("Accept", "text/plain").
		End()

	if errs != nil {
		err := errs[len(errs)-1]
		log.Fatal(err.Error())
	}

	return string(body)
}

func noErrorOn404(err error) error {

	if err != nil {
		errStatusCode := err.(misc.HTTPStatus)
		if errStatusCode.StatusCode != 404 {
			// fmt.Println("error found in noerroron404", err)
			return err
		}
	}
	return nil
}

// PortByName - sort functions for network ports by name
type PortByName []network.PortResponse

func (a PortByName) Len() int           { return len(a) }
func (a PortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PortByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// ServerByName - sort functions for server by name
type ServerByName []compute.Server

func (a ServerByName) Len() int           { return len(a) }
func (a ServerByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ServerByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
