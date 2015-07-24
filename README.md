#hpcloud-kubesetup 
##Deploying Kubernetes clusters to HP Helion OpenStack
##====================================================

This repository contains the code and instructions for the hpcloud-kubesetup installer tool. The hpcloud-kubesetup installer enables you to create and deploy Kubernetes (1.0.1) clusters on to your own private HP Helion OpenStack environment (version 1.1 or later) or to your hosted HP Helion Public Cloud account. 

The installer process runs on your workstation, provisioning the cluster remotely.

## Prerequisites ##
1. Credentials to your HP Helion OpenStack environment or HP Helion Public Cloud account.
2. CoreOS version 653.0.0 or later loaded in to OpenStack glance  [(steps)](https://coreos.com/os/docs/latest/booting-on-openstack.html). Note: when deploying to a HP Helion Public Cloud account this prerquisite is already satisfied.
3. An OpenStack project/tenant to deploy your Kubernetes cluster to. Note: when deploying to a HP Helion Public Cloud account, you can use the existing tenant.
3. A private network within the OpenStack project/tenant, providingt network isolation [(steps)](https://github.com/gertd/hpcloud-kubesetup/blob/master/scripts/create-private-network.sh).
4. Ingress TCP communication over ports 22 (SSH), 80, 443 and 8080 (kube-apiserver) by adding these rulese to the default OpenStack security group within the project [(steps)](https://github.com/gertd/hpcloud-kubesetup/blob/master/scripts/update-default-securitygroup.sh)
5. A Linux, Mac, or Windows workstation with internet connectivity and connectivity to your HP Helion OpenStack environment.

## Steps ##
1. Download and install the hpcloud-kubesetup installer and Kubernetes kubectl utility for your specific platform:

	**Linux**

		mkdir -p ~/kube
		
		wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-linux.zip \
		-O ~/kube/hpcloud-kubesetup-linux.zip
		unzip hpcloud-kubesetup-linux.zip -d ~/kube/
		sudo mv ~/kube/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup
		sudo chmod +x /usr/local/bin/hpcloud-kubesetup
		
		wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl \
		-O ~/kube/kubectl
		sudo mv ~/kube/kubectl /usr/local/bin/kubectl
		sudo chmod +x /usr/local/bin/kubectl
	
	**Mac**

		mkdir -p ~/kube
		
		wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-darwin.zip \
		-O ~/kube/hpcloud-kubesetup-darwin.zip
		unzip hpcloud-kubesetup-darwin.zip -d ~/kube/
		sudo mv ~/kube/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup
		sudo chmod +x /usr/local/bin/hpcloud-kubesetup
		
		wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/darwin/amd64/kubectl \
		-O ~/kube/kubectl
		sudo mv ~/kube/kubectl /usr/local/bin/kubectl
		sudo chmod +x /usr/local/bin/kubectl
		
	**Windows**
	
	[Installation script](https://github.com/gertd/hpcloud-kubesetup/blob/master/windows/README.md)
	For manual installation steps:
	1. Download [hpcloud-kubesetup-windows.zip](https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-windows.zip) 
	2. Unzip hpcloud-kubesetup.zip
	3. Download [kubectl.exe](https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/windows/amd64/kubectl.exe) 
	
2. Log into your account and download the "OpenStack RC file" located on the Project\Access & Security panel inside the API Access tab. The [download button](https://a248.e.akamai.net/cdn.hpcloudsvc.com/ha4ca03ecf0c27c00f0c991360b263f06/prodaw2/rc-file.png) is on the top right corner.

3. Set environment variables 

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
	
	Open the OpenStack resource script with a text editor such as Notepad++. Replace the variables with your configuration values. Then run this in your command prompt window.
	
		set OS_AUTH_URL=<OS_AUTH_URL>
		set OS_TENANT_ID=<OS_TENANT_ID>
		set OS_TENANT_NAME=<OS_TENANT_NAME>
		set OS_USERNAME=<OS_USERNAME>
		set OS_PASSWORD=<OS_PASSWORD>	
		set OS_REGION_NAME=<OS_REGION_NAME>

4. Update `kubesetup.yml` if necessary. This file describes the setup of the cluster. By default, a cluster consisting of 3 nodes, 1 master node and 2 minion nodes, will be created. 

	You will need to:
	 * Create a new ssh key named `kube-key` or modify `sshkey` to reflect the key you want to use instead
	 * Create or modify the network
	 * Verify and update if needed the ip based on the ip range on your tenant

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
	>hpcloud-kubesetup.exe install
	2015/07/23 12:06:23 config file          - kube-master {10.0.0.140 true CoreOS standard.medium }
	2015/07/23 12:06:23 config file          - kube-node-1 {10.0.0.141 false CoreOS standard.small }
	2015/07/23 12:06:23 config file          - kube-node-2 {10.0.0.142 false CoreOS standard.small }
	2015/07/23 12:06:23 config file          - SSHKey <redacted>
	2015/07/23 12:06:23 config file          - Network CloudHorizonNetwork
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
	2015/07/23 12:06:25 create port          - kube-master 10.0.0.140
	2015/07/23 12:06:26 create port          - 86587b0b-4351-467e-baaf-882a6f71f952 COMPLETED
	2015/07/23 12:06:26 create server        - kube-master 10.0.0.140
	2015/07/23 12:06:27 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:27 flavor               - 102
	2015/07/23 12:06:28 create server        - password pAYv4g9bH4Wj
	2015/07/23 12:06:28 create server        - a77e155a-847f-41b8-a523-6d14a044a568 COMPLETED
	2015/07/23 12:06:28 create port          - kube-node-1 10.0.0.141
	2015/07/23 12:06:28 create port          - fb1180ea-134d-477f-a9a0-ad1e1ea9e447 COMPLETED
	2015/07/23 12:06:28 create server        - kube-node-1 10.0.0.141
	2015/07/23 12:06:29 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:29 flavor               - 101
	2015/07/23 12:06:29 create server        - password CynjWU23GgBb
	2015/07/23 12:06:29 create server        - 2627034a-6673-4837-976f-2620f4e4af4a COMPLETED
	2015/07/23 12:06:29 create port          - kube-node-2 10.0.0.142
	2015/07/23 12:06:30 create port          - a9a62294-9ce8-4804-8a93-3f0d5808b19a COMPLETED
	2015/07/23 12:06:30 create server        - kube-node-2 10.0.0.142
	2015/07/23 12:06:30 image                - 5c2ccd59-1ae8-417a-8abc-22fb4f4b9f85
	2015/07/23 12:06:30 flavor               - 101
	2015/07/23 12:06:31 create server        - password BeyKKM3WecAg
	2015/07/23 12:06:31 create server        - 5bae49a3-e1c4-4a3e-8443-31702442a4e7 COMPLETED
	2015/07/23 12:06:31 server status        - kube-master BUILD
	2015/07/23 12:06:54 server status        - kube-master ACTIVE
	2015/07/23 12:06:54 server status        - kube-node-1 ACTIVE
	2015/07/23 12:06:54 server status        - kube-node-2 ACTIVE
	2015/07/23 12:06:54 associate IP         - kube-master 15.125.106.149
	2015/07/23 12:06:55 associate IP         - kube-master COMPLETED	
	```

7. The installer associates a floating IP address with the Kubernetes master node. You can get the floating IP in your instances Horizon panel. The next step is to ssh in to the master node and run the Kubernetes kubecfg tool to list the minions and verify everything is working properly.

	**Linux**

		ssh -i kube-key core@15.126.200.248

		kubecfg get nodess 
		
		Results:
	
			$ ssh -i ../kube-key core@15.126.200.248
			
			The authenticity of host '15.126.200.248 (15.126.200.248)' can't be established.
			RSA key fingerprint is fe:b1:a0:6f:3b:60:e7:3c:26:30:98:4a:86:24:99:d8.
			Are you sure you want to continue connecting (yes/no)? yes
			Warning: Permanently added '15.126.200.248' (RSA) to the list of known hosts.
			CoreOS (stable)
			core@kube-master ~ $ kubecfg get nodes
			NAME         LABELS                              STATUS
			10.0.0.141   kubernetes.io/hostname=10.0.0.141   Ready
			10.0.0.142   kubernetes.io/hostname=10.0.0.142   Ready
	
	**MacOS**

		ssh -i kube-key core@15.126.200.248
			
		kubecfg get nodes
		
		Results:

			$ ssh -i ../kube-key core@15.126.200.248
			
			The authenticity of host '15.126.200.248 (15.126.200.248)' can't be established.
			RSA key fingerprint is fe:b1:a0:6f:3b:60:e7:3c:26:30:98:4a:86:24:99:d8.
			Are you sure you want to continue connecting (yes/no)? yes
			Warning: Permanently added '15.126.200.248' (RSA) to the list of known hosts.
			CoreOS (stable)
			core@kube-master ~ $ kubecfg get nodes
			NAME         LABELS                              STATUS
			10.0.0.141   kubernetes.io/hostname=10.0.0.141   Ready
			10.0.0.142   kubernetes.io/hostname=10.0.0.142   Ready
	
	**Windows**
	
	You will need an ssh client such as cywin for this.

		>ssh -i my_key.pem core@15.126.200.248 "/opt/bin/kubectl get nodes"

		Results:

			>ssh -i my_key.pem core@15.126.200.248 "/opt/bin/kubectl get nodes"
			Warning: Permanently added '15.126.200.248' (ED25519) to the list of known hosts.
			NAME         LABELS                              STATUS
			10.0.0.141   kubernetes.io/hostname=10.0.0.141   Ready
			10.0.0.142   kubernetes.io/hostname=10.0.0.142   Ready


8. After verifying all the nodes are there, you are ready to rock and roll. The next step would be to deploy a [sample application](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/examples/guestbook/README.md
) to the cluster. Happy containerizing!

## License ##

Copyright 2015 Hewlett-Packard

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
