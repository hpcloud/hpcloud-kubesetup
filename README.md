#hpcloud-kubesetup
##Deploying Kubernetes clusters to HP Helion OpenStack

This repository contains the code and instructions for the hpcloud-kubesetup installer tool. The hpcloud-kubesetup installer enables you to create and deploy Kubernetes (1.0.1) clusters on to your own private HP Helion OpenStack environment (version 1.1 or later) or to your hosted HP Helion Public Cloud account.

The installer process runs on your workstation, provisioning the cluster remotely.

## Prerequisites ##
1. Credentials to your HP Helion OpenStack environment or HP Helion Public Cloud account.
2. CoreOS version 653.0.0 or later loaded in to OpenStack glance  [(steps)](https://coreos.com/os/docs/latest/booting-on-openstack.html). Note: when deploying to a HP Helion Public Cloud account this prerquisite is already satisfied.
3. An OpenStack project/tenant to deploy your Kubernetes cluster to. Note: when deploying to a HP Helion Public Cloud account, you can use the existing tenant.
3. A private network within the OpenStack project/tenant, providing network isolation [(steps)](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/scripts/create-private-network.sh).
4. Ingress TCP communication over ports 22 (SSH), 80, 443 and 8080 (kube-apiserver) by adding these rulese to the default OpenStack security group within the project [(steps)](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/scripts/update-default-securitygroup.sh)
5. A Linux, Mac, or Windows workstation with internet connectivity and connectivity to your HP Helion OpenStack environment.

## Steps ##
1. Download and install the hpcloud-kubesetup installer and Kubernetes kubectl utility for your specific platform:

	**Linux**

	Script based installation folllow these [instructions](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/setup/linux/README.md).

	Manual installation steps:

		mkdir -p /usr/local/kubernetes

		wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-linux.zip \
		-O /usr/local/kubernetes/hpcloud-kubesetup-linux.zip

		unzip -o /usr/local/kubernetes/hpcloud-kubesetup-linux.zip -d /usr/local/kubernetes/

		wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl \
		-O /usr/local/kubernetes/linux/kubectl

		chmod +x /usr/local/kubernetes/linux/hpcloud-kubesetup
		ln -s /usr/local/kubernetes/linux/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup

		chmod +x /usr/local/kubernetes/linux/kubectl
		ln -s /usr/local/kubernetes/linux/kubectl /usr/local/bin/kubectl

		mkdir -p ~/kubernetes
		cp -n /usr/local/kubernetes/linux/kubesetup.yml ~/kubernetes/.

	**Mac**

	Script based installation folllow these [instructions](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/setup/darwin/README.md).

	Manual installation steps:

		mkdir -p /usr/local/kubernetes

		wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-darwin.zip \
		-O /usr/local/kubernetes/hpcloud-kubesetup-darwin.zip

		unzip -o /usr/local/kubernetes/hpcloud-kubesetup-darwin.zip -d /usr/local/kubernetes/

		wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/darwin/amd64/kubectl \
		-O /usr/local/kubernetes/darwin/kubectl

		chmod +x /usr/local/kubernetes/darwin/hpcloud-kubesetup
		ln -s /usr/local/kubernetes/darwin/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup

		chmod +x /usr/local/kubernetes/darwin/kubectl
		ln -s /usr/local/kubernetes/darwin/kubectl /usr/local/bin/kubectl

		mkdir -p ~/kubernetes
		cp -n /usr/local/kubernetes/darwin/kubesetup.yml ~/kubernetes/.

	**Windows**

	Script based installation folllow these [instructions](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/setup/windows/README.md).

	Manual installation steps:

	1. Download [hpcloud-kubesetup-windows.zip](https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-windows.zip)
	2. Unzip hpcloud-kubesetup.zip
	3. Download [kubectl.exe](https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/windows/amd64/kubectl.exe)

2. Log into the OpenStack Horizon portal with your account and download the "OpenStack RC file" located on the Project\Access & Security panel inside the API Access tab. The [download button](https://a248.e.akamai.net/cdn.hpcloudsvc.com/ha4ca03ecf0c27c00f0c991360b263f06/prodaw2/rc-file.png) is on the top right corner.

3. Setup OpenStack environment variables

	**Mac & Linux**

	Execute the OpenStack resource script. The script will ask you to enter your OpenStack password. All settings will be exported as environment variables.

		source ./<your project name>-openrc.sh

	To inspect what was exported, run `export | grep OS_`. You should see a similar result to:

		$ export | grep OS_
		declare -x OS_AUTH_URL="https://region-a.geo-1.identity.hpcloudsvc.com:35357/v2.0/"
		declare -x OS_PASSWORD="My Very Secret Password"
		declare -x OS_TENANT_ID="12345678901234"
		declare -x OS_TENANT_NAME="kubernetes"
		declare -x OS_USERNAME="kube"

	**Windows**

	Rename the downloaded <your project name>-openrc.sh file to <your project name>-openrc.bat
	Open the <your project name>-openrc.bat file within an editor like notepad.
	Replace the export statement with set statement, like shown below.

		set OS_AUTH_URL=<OS_AUTH_URL>
		set OS_TENANT_ID=<OS_TENANT_ID>
		set OS_TENANT_NAME=<OS_TENANT_NAME>
		set OS_USERNAME=<OS_USERNAME>
		set OS_PASSWORD=<OS_PASSWORD>
		set OS_REGION_NAME=<OS_REGION_NAME>

	Run the <your project name>-openrc.bat file inside the console window from which we will the remaining installer steops

4. Update `kubesetup.yml` if necessary. This file describes the setup of the cluster. By default, a cluster consisting of 3 nodes, 1 master node and 2 minion nodes, will be created.

	You will need to:
	 * Create a new ssh key named `kube-key` or modify `sshkey` to reflect the key name of an existing key pair inside OpenStack
	 * Create the kube-net network [(steps)](https://github.com/hpcloud/hpcloud-kubesetup/blob/master/scripts/create-private-network.sh) or modify the network entry in the kubesetup.yml file to an existing private network inside the project/tenant you will be deploying to
	 * Verify if specified IP address range is supported by your subnet. When using the create-private-network.sh script you can use the default values

	**kubesetup.yml**

		hosts:
		  kube-master:
		    ip: 192.168.1.140
		    ismaster: true
		    vm-image: CoreOS
		    vm-size: standard.medium
		  kube-node-1:
		    ip: 192.168.1.141
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small
		  kube-node-2:
		    ip: 192.168.1.142
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small

		sshkey: kube-key
		network: kube-net
		availabilityZone: az2

6. Once your `kubesetup.yml` reflects the type of cluster you want to create, you can then execute the cluster installer:

	**Mac & Linux**

		hpcloud-kubesetup install

	**Windows**

		hpcloud-kubesetup.exe install


	Once run, you should see the following results:

	```
	$ hpcloud-kubesetup install
	2015/07/23 12:06:23 config file          - kube-master {192.168.1.140 true CoreOS standard.medium }
	2015/07/23 12:06:23 config file          - kube-node-1 {192.168.1.141 false CoreOS standard.small }
	2015/07/23 12:06:23 config file          - kube-node-2 {192.168.1.142 false CoreOS standard.small }
	2015/07/23 12:06:23 config file          - SSHKey <redacted>
	2015/07/23 12:06:23 config file          - Network kube-net
	2015/07/23 12:06:23 config file          - AvailabilityZone az2
	2015/07/23 12:06:23 OS_AUTH_URL          -  <redacted>
	2015/07/23 12:06:23 OS_TENANT_ID         -  <redacted>
	2015/07/23 12:06:23 OS_TENANT_NAME       -  <redacted>
	2015/07/23 12:06:23 OS_USERNAME          -  <redacted>
	2015/07/23 12:06:23 OS_REGION_NAME       -  <redacted>
	2015/07/23 12:06:23 OS_AUTH_TOKEN        -
	2015/07/23 12:06:23 OS_CACERT            -
	2015/07/23 12:06:23 skip-ssl-validation  - false
	2015/07/23 12:06:23 debug                - false
	2015/07/23 12:06:24 token                - HPAuth10_9b0328c27a31d4c4ff52cbd447270a8bc909572cbccf76b770c1cb06cc9f1986
	2015/07/23 12:06:25 network              - aca348f6-b481-469b-8aef-efd235987578
	2015/07/23 12:06:25 create cloudconfig   - kube-master
	2015/07/23 12:06:25 create cloudconfig   - kube-master.yml COMPLETED
	2015/07/23 12:06:25 create cloudconfig   - kube-node-1
	2015/07/23 12:06:25 create cloudconfig   - kube-node-1.yml COMPLETED
	2015/07/23 12:06:25 create cloudconfig   - kube-node-2
	2015/07/23 12:06:25 create cloudconfig   - kube-node-2.yml COMPLETED
	2015/07/23 12:06:25 create port          - kube-master 192.1.168.140
	2015/07/23 12:06:26 create port          - 86587b0b-4351-467e-baaf-882a6f71f952 COMPLETED
	2015/07/23 12:06:26 create server        - kube-master 192.168.1.140
	2015/07/23 12:06:27 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:27 flavor               - 102
	2015/07/23 12:06:28 create server        - password <redacted>
	2015/07/23 12:06:28 create server        - a77e155a-847f-41b8-a523-6d14a044a568 COMPLETED
	2015/07/23 12:06:28 create port          - kube-node-1 192.168.1.141
	2015/07/23 12:06:28 create port          - fb1180ea-134d-477f-a9a0-ad1e1ea9e447 COMPLETED
	2015/07/23 12:06:28 create server        - kube-node-1 192.168.1.141
	2015/07/23 12:06:29 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:29 flavor               - 101
	2015/07/23 12:06:29 create server        - password <redacted>
	2015/07/23 12:06:29 create server        - 2627034a-6673-4837-976f-2620f4e4af4a COMPLETED
	2015/07/23 12:06:29 create port          - kube-node-2 192.168.1.142
	2015/07/23 12:06:30 create port          - a9a62294-9ce8-4804-8a93-3f0d5808b19a COMPLETED
	2015/07/23 12:06:30 create server        - kube-node-2 192.168.1.142
	2015/07/23 12:06:30 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:30 flavor               - 101
	2015/07/23 12:06:31 create server        - password <redacted>
	2015/07/23 12:06:31 create server        - 5bae49a3-e1c4-4a3e-8443-31702442a4e7 COMPLETED
	2015/07/23 12:06:31 server status        - kube-master BUILD
	2015/07/23 12:06:54 server status        - kube-master ACTIVE
	2015/07/23 12:06:54 server status        - kube-node-1 ACTIVE
	2015/07/23 12:06:54 server status        - kube-node-2 ACTIVE
	2015/07/23 12:06:54 associate IP         - kube-master 15.125.106.149
	2015/07/23 12:06:55 associate IP         - kube-master COMPLETED
	```

7. The installer associates a floating IP address with the Kubernetes master node. You can find the floating IP in list if server instances in the Horizon panel or by using the nova list command. The next step is use kubectl to explore and inspect the cluster.

	**Mac & Linux & Windows**

		$ kubectl cluster-info --server=http://15.125.106.149:8080
		Kubernetes master is running at http://15.125.106.149:8080

		$ kubectl version --server=http://15.125.106.149:8080
		Client Version: version.Info{Major:"1", Minor:"0", GitVersion:"v1.0.1", GitCommit:"6a5c06e3d1eb27a6310a09270e4a5fb1afa93e74", GitTreeState:"clean"}
		Server Version: version.Info{Major:"1", Minor:"0", GitVersion:"v1.0.1", GitCommit:"6a5c06e3d1eb27a6310a09270e4a5fb1afa93e74", GitTreeState:"clean"}

		$ kubectl get nodes --server=http://15.125.106.149:8080
		NAME            LABELS                                 STATUS
		192.168.1.141   kubernetes.io/hostname=192.168.1.141   Ready
		192.168.1.142   kubernetes.io/hostname=192.168.1.142   Ready
		

	Alternatively for Mac & Linux you can setup a secure SSH tunel between the kubectl client and the kube-apiserver, this prevents from having to provide the --server parameter on each call. The confige the SSH tunel use the following command:
	
		ssh -f -nNT -L 8080:127.0.0.1:8080 core@<master-public-ip>
		
		ssh -f -nNT -L 8080:127.0.0.1:8080 core@15.125.106.149
		$ kubectl get services
		NAME         LABELS                                    SELECTOR   IP(S)        PORT(S)
		kubernetes   component=apiserver,provider=kubernetes   <none>     10.100.0.1   443/TCP

8. After verifying all the nodes are there, you are ready to rock and roll. The next step will be to deploy a [sample application](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/examples/guestbook/README.md
) to your Kubernetes cluster!

Happy containerizing!

## License ##

Copyright 2015 Hewlett-Packard

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
