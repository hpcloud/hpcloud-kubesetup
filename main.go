// hpcloud-kubesetup project main.go
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gertd/go-openstack"
	"github.com/gertd/go-openstack/compute"
	"github.com/gertd/go-openstack/identity"
	"github.com/gertd/go-openstack/image"
	"github.com/gertd/go-openstack/network"

	"github.com/parnurzeal/gorequest"

	"gopkg.in/yaml.v1"
)

var (
	configFile string
	authUrl    string
	tenantId   string
	tenantName string
	username   string
	password   string
	region     string
	uninstall  bool
)

type Config struct {
	Nodes  map[string]Node `yaml:"hosts"`
	SSHKey string          `yaml:"sshkey"`
}

type Node struct {
	IP       string `yaml:"ip"`
	IsMaster bool   `yaml:"ismaster"`
	VMImage  string `yaml:"vm-image"`
	VMSize   string `yaml:"vm-size"`
	ServerId string
}

func init() {
	flag.StringVar(&configFile, "c", "kubesetup.yml", "Kubernetes cluster configuration file")
	flag.StringVar(&authUrl, "a", "", "OpenStack authentication URL (OS_AUTH_URL)")
	flag.StringVar(&tenantId, "i", "", "OpenStack tenant id (OS_TENANT_ID)")
	flag.StringVar(&tenantName, "n", "", "OpenStack tenant name (OS_TENANT_NAME)")
	flag.StringVar(&username, "u", "", "OpenStack user name (OS_USERNAME)")
	flag.StringVar(&password, "p", "", "OpenStack passsword (OS_PASSWORD)")
	flag.StringVar(&region, "r", "", "OpenStack region name (OS_REGION_NAME)")
	flag.BoolVar(&uninstall, "U", false, "Uninstall cluster (defined in config file -c)")
}

func main() {
	flag.Parse()

	config, err := readConfigFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.Log()

	openStackConfig, err := openstack.InitializeFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}
	updateConfigFromCommandLine(&openStackConfig)
	openStackConfig.Log()

	auth, err := identity.Authenticate(openStackConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%-20s - %s\n", "token:", auth.Access.Token.Id)

	keypair, err := compute.GetKeypair(auth, config.SSHKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	subnets, err := network.GetSubnets(auth)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	ports, err := network.GetPorts(auth)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}
	sort.Sort(network.ByName(ports))

	servers, err := compute.GetServers(auth)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}
	sort.Sort(compute.ByName(servers))

	// uninstall/cleanup
	for _, v := range servers {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete server", v.Name)

			err := compute.DeleteServer(auth, v.Id)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete server", v.Name, "COMPLETED")
		}
	}

	for _, v := range ports {

		if _, ok := config.Nodes[v.Name]; ok {

			log.Printf("%-20s - %s\n", "delete port", v.Name)

			err := network.DeletePort(auth, v.Id)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}

			log.Printf("%-20s - %s %s\n", "delete port", v.Name, "COMPLETED")
		}
	}

	if uninstall {
		os.Exit(1)
	}

	// create cloudconfig files
	nodeId := 0
	for k, v := range config.Nodes {

		log.Printf("%-20s - %s\n", "create cloudconfig", k)

		data := make(map[string]string)

		data["filename"] = k + ".yml"

		if v.IsMaster {
			data["role"] = "master"
		} else {
			data["role"] = "minion"
		}

		data["machines"] = getMachines(config.Nodes)
		data["peers"] = getPeers(config.Nodes, v)
		data["hostname"] = k
		data["ip"] = v.IP
		data["dns"] = subnets[0].GatewayIP
		data["gateway"] = subnets[0].GatewayIP
		data["subnet"] = fmt.Sprintf("10.244.%d.1/24", nodeId) // subnets.Subnets[0].CIDR
		data["sshkey"] = keypair.PublicKey
		data["discovery"] = getDiscoveryKey()

		nodeId = nodeId + 1

		createCloudConfig(data)

		log.Printf("%-20s - %s %s\n", "create cloudconfig", data["filename"], "COMPLETED")
	}

	// install
	for k, v := range config.Nodes {

		log.Printf("%-20s - %s %s\n", "create port", k, v.IP)

		newPort := network.Port{}
		newPort.Name = k
		newPort.NetworkId = subnets[0].NetworkId
		newPort.FixedIPs = []network.FixedIP{{subnets[0].Id, v.IP}}

		port, err := network.CreatePort(auth, newPort)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		log.Printf("%-20s - %s %s\n", "create port", port.Id, "COMPLETED")

		log.Printf("%-20s - %s %s\n", "create server", k, v.IP)

		userdata, err := getUserData(k + ".yml")
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		flavor, err := compute.GetFlavor(auth, v.VMSize)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}
		log.Printf("%-20s - %s\n", "flavor:", flavor.Id)

		image, err := image.GetImage(auth, v.VMImage)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}
		log.Printf("%-20s - %s\n", "image:", image.Id)

		newServer := compute.NewServer{}
		newServer.Name = k
		newServer.ImageRef = image.Id
		newServer.FlavorRef = flavor.Id
		newServer.KeyName = keypair.Name
		newServer.UserData = userdata
		newServer.Network = []compute.Network{{port.NetworkId, port.Id}}
		newServer.SecurityGroups = []compute.SecurityGroup{{"default"}}

		server, err := compute.CreateServer(auth, newServer)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		log.Printf("%-20s - %s %s\n", "create server", "password", server.AdminPass)
		log.Printf("%-20s - %s %s\n", "create server", server.Id, "COMPLETED")

		node := config.Nodes[k]
		node.ServerId = server.Id
		config.Nodes[k] = node
	}

	// status
	for _, v := range config.Nodes {

		prevStatus := ""
		for {
			server, err := compute.GetServer(auth, v.ServerId)
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

	// associate floating IP with each master node
	var unAssigned []compute.FloatingIP
	floatingIPs, err := compute.GetFloatingIPs(auth)
	if err != nil {
		log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
	}

	for _, v := range floatingIPs {
		if len(v.InstanceId) == 0 {
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
			fp, err := compute.CreateFloatingIP(auth)
			if err != nil {
				log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
			}
			unAssigned = append(unAssigned, fp)

			log.Printf("%-20s - %s %s\n", "create public IP", fp.IP, "COMPLETED")

		}

		log.Printf("%-20s - %s %s\n", "associate IP", k, unAssigned[0].IP)

		err := compute.ServerAction(auth, v.ServerId, "addFloatingIp", "address", unAssigned[0].IP)
		if err != nil {
			log.Fatal(fmt.Sprintf("%-20s - %s\n", "error:", err.Error()))
		}

		unAssigned = append(unAssigned[1:])
		log.Printf("%-20s - %s %s\n", "associate IP", k, "COMPLETED")
	}

	os.Exit(1)
}

func updateConfigFromCommandLine(config *openstack.OpenStackConfig) {

	if len(authUrl) > 0 {
		config.AuthUrl = authUrl
	}
	if len(tenantId) > 0 {
		config.TenantId = tenantId
	}
	if len(tenantName) > 0 {
		config.TenantName = tenantName
	}
	if len(username) > 0 {
		config.Username = username
	}
	if len(password) > 0 {
		config.Password = password
	}
	if len(region) > 0 {
		config.Region = region
	}
}

func readConfigFile(filename string) (config Config, err error) {

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

func (config Config) Log() {
	for k, v := range config.Nodes {
		log.Printf("%-20s - %s %v\n", "config file", k, v)
	}
	log.Printf("%-20s - %s %s\n", "config file", "SSHKey", config.SSHKey)
}

func createCloudConfig(data map[string]string) error {

	var b bytes.Buffer
	f, err := os.Create(data["filename"])
	if err != nil {
		return err
	}

	w := io.MultiWriter(f, &b)
	if err := nodeTmpl.Execute(w, data); err != nil {
		return err
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

func getMachines(nodeList map[string]Node) string {

	var csv bytes.Buffer

	i := 0
	for _, v := range nodeList {
		csv.WriteString(v.IP)
		if i < (len(nodeList) - 1) {
			csv.WriteString(",")
		}
		i = i + 1
	}

	return csv.String()
}

func getPeers(nodeList map[string]Node, self Node) string {

	var csv bytes.Buffer

	i := 0
	for _, v := range nodeList {

		if v.IP == self.IP {
			continue
		}

		csv.WriteString(v.IP + ":7001")
		if i < (len(nodeList) - 2) {
			csv.WriteString(",")
		}
		i = i + 1
	}

	return csv.String()
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
