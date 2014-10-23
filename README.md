hpcloud-kubesetup
=================

This repository contains the code and instructions for the hpcloud-kubesetup installer, which deploys a Kubernetes cluster to the HP Helion Public Cloud or to any HP Helion OpenStack environment. The process runs on your workstation, provisioning the cluster remotely.

## Prerequisites ##
1. A valid hpcloud.com account or valid credentials to your Helion OpenStack environment.
2. CoreOS available as a guest image if using your own Helion OpenStack environment.
3. A Linux, Mac, or Windows workstation with internet connectivity.

## Steps ##
1. Download the hpcloud-kubesetup installer for your specific platform:

	**Linux**

		wget https://github.com/hpcloud/hpcloud-kubesetup/bin/linux/hpcloud-kubesetup_linux_amd64 -O /usr/local/bin/hpcloud-kubesetup
		chmod +x /usr/local/bin/hpcloud-kubesetup
	
	**Mac**

		wget https://github.com/hpcloud/hpcloud-kubesetup/bin/darwin/hpcloud-kubesetup_darwin_amd64 -O /usr/local/bin/hpcloud-kubesetup
		chmod +x /usr/local/bin/hpcloud-kubesetup
		
	**Windows**

		wget https://github.com/hpcloud/hpcloud-kubesetup/bin/windows/hpcloud-kubesetup_windows_amd64.exe -O %temp%/hpcloud-kubesetup.exe
	
2. Download the default setup configuration file to create a 4 node Kubernetes cluster

	**Linux, Mac, Windows**

		wget https://github.com/hpcloud/kubesetup/blob/master/kubesetup.yml 
		

3. Log into your account and download the "OpenStack RC file" located on the Project\Access & Security panel inside the API Access tab. The [download button](https://a248.e.akamai.net/cdn.hpcloudsvc.com/ha4ca03ecf0c27c00f0c991360b263f06/prodaw2/rc-file.png) is on the top right corner.

4. Execute the OpenStack resource script. The script will ask you to enter your OpenStack password. All settings will be exported as environment variables.

	**Linux**

		source ./<your project name>-openrc.sh
	
	**Mac**
	
		source ./<your project name>-openrc.sh
		
	**Windows**

		Coming soon

	To inspect what was exported, run `export | grep OS_`. You should see a similar result to:
		
		$ export | grep OS_
		declare -x OS_AUTH_URL="https://region-a.geo-1.identity.hpcloudsvc.com:35357/v2.0/"
		declare -x OS_PASSWORD="My Very Secret Password"
		declare -x OS_TENANT_ID="12345678901234"
		declare -x OS_TENANT_NAME="kubernetes"
		declare -x OS_USERNAME="kube"

5. Update `kubesetup.yml` if needed. This file describes the setup of the cluster. By default, a cluster consisting of 5 nodes, 1 master node and 4 minion nodes, will be created. You will need to create a new ssh key named `kube-key` or modify `sshkey` to reflect the key you want to use instead.

	*kubesetup.yml*

	    hosts:
	      kube-master:
	    	ip: 10.0.0.40
	    	ismaster: true
	    	vm-image: CoreOS
	    	vm-size: standard.medium
	      kube-minion-1:
	    	ip: 10.0.0.41
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small
	      kube-minion-2:
		    ip: 10.0.0.42
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small
	      kube-minion-3:
		    ip: 10.0.0.43
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small
	      kube-minion-4:
		    ip: 10.0.0.44
		    ismaster: false
		    vm-image: CoreOS
		    vm-size: standard.small

		sshkey: kube-key

6. Once `kubesetup.yml` reflects the type of cluster you want to create, you can then execute the cluster installer:

	**Linux**

		/usr/local/bin/hpcloud-kubesetup

	**Mac**
		
		./usr/local/bin/hpcloud-kubesetup
	
	Windows

		coming soon


	Once run, you should see the following results:
	
	
		$ ./hpcloud-kubesetup
		2014/10/13 14:13:56 config file          - kube-master {10.0.0.40 true CoreOS standard.medium }
		2014/10/13 14:13:56 config file          - kube-minion-1 {10.0.0.41 false CoreOS standard.small }
		2014/10/13 14:13:56 config file          - kube-minion-2 {10.0.0.42 false CoreOS standard.small }
		2014/10/13 14:13:56 config file          - SSHKey kube-key
		2014/10/13 14:13:56 OS_AUTH_URL          - https://region-a.geo-1.identity.hpcloudsvc.com:35357/v2.0/
		2014/10/13 14:13:56 OS_TENANT_ID         - 12345678901234
		2014/10/13 14:13:56 OS_TENANT_NAME       - kubernetes
		2014/10/13 14:13:56 OS_USERNAME          - kube
		2014/10/13 14:13:56 OS_REGION_NAME       -
		2014/10/13 14:13:57 token:               - HPAuth10_12d08ac1f1f09dad8ae89dfc19b82d4b554996458858a9ae0b579829535cfed6
		2014/10/13 14:13:59 create cloudconfig   - kube-master
		2014/10/13 14:14:00 create cloudconfig   - kube-master.yml COMPLETED
		2014/10/13 14:14:00 create cloudconfig   - kube-minion-1
		2014/10/13 14:14:00 create cloudconfig   - kube-minion-1.yml COMPLETED
		2014/10/13 14:14:00 create cloudconfig   - kube-minion-2
		2014/10/13 14:14:00 create cloudconfig   - kube-minion-2.yml COMPLETED
		2014/10/13 14:14:00 create port          - kube-master 10.0.0.40
		2014/10/13 14:14:01 create port          - ae517467-ab85-4498-ae30-7afd9dc8a4d6 COMPLETED
		2014/10/13 14:14:01 create server        - kube-master 10.0.0.40
		2014/10/13 14:14:02 flavor:              - 102
		2014/10/13 14:14:03 image:               - e53e5521-3f42-455f-bd9d-16cd4d5776e5
		2014/10/13 14:14:04 create server        - password B9p7AFCCgvqH
		2014/10/13 14:14:04 create server        - 99112dec-4d45-40e9-a87c-79a45276f041 COMPLETED
		2014/10/13 14:14:04 create port          - kube-minion-1 10.0.0.41
		2014/10/13 14:14:05 create port          - 8b218e77-db46-419e-b0c4-59b1e2d1fbc8 COMPLETED
		2014/10/13 14:14:05 create server        - kube-minion-1 10.0.0.41
		2014/10/13 14:14:05 flavor:              - 101
		2014/10/13 14:14:06 image:               - e53e5521-3f42-455f-bd9d-16cd4d5776e5
		2014/10/13 14:14:07 create server        - password 5Ba7Z9pydpUG
		2014/10/13 14:14:07 create server        - b9c08634-1cf9-450a-9acf-bbf838bc3404 COMPLETED
		2014/10/13 14:14:07 create port          - kube-minion-2 10.0.0.42
		2014/10/13 14:14:08 create port          - b0c543a0-e4bf-41bc-b4b7-b0bd4d918e56 COMPLETED
		2014/10/13 14:14:08 create server        - kube-minion-2 10.0.0.42
		2014/10/13 14:14:09 flavor:              - 101
		2014/10/13 14:14:09 image:               - e53e5521-3f42-455f-bd9d-16cd4d5776e5
		2014/10/13 14:14:11 create server        - password QTFAJsk5Ehde
		2014/10/13 14:14:11 create server        - c4282263-950e-41fa-ad1d-1343a27db2c2 COMPLETED
		2014/10/13 14:14:11 server status        - kube-master BUILD
		2014/10/13 14:14:30 server status        - kube-master ACTIVE
		2014/10/13 14:14:31 server status        - kube-minion-1 ACTIVE
		2014/10/13 14:14:32 server status        - kube-minion-2 ACTIVE
		2014/10/13 14:14:32 associate IP         - kube-master 15.126.200.248
		2014/10/13 14:14:34 associate IP         - kube-master COMPLETED
			
	
7. The installer associates a floating IP address with the Kubernetes master node. The next step is to ssh in to the master node and run the Kubernetes kubecfg tool to list the minions and verify everything is working properly.

	**Linux**

		ssh -i kube-key core@15.126.200.248

		kubecfg list minions 
		
		Results:
	
			$ ssh -i ../kube-key core@15.126.200.248
			
			The authenticity of host '15.126.200.248 (15.126.200.248)' can't be established.
			RSA key fingerprint is fe:b1:a0:6f:3b:60:e7:3c:26:30:98:4a:86:24:99:d8.
			Are you sure you want to continue connecting (yes/no)? yes
			Warning: Permanently added '15.126.200.248' (RSA) to the list of known hosts.
			CoreOS (stable)
			core@kube-master ~ $ kubecfg list minions
			Minion identifier
			----------
			10.0.0.40
			10.0.0.41
			10.0.0.42
	
	**MacOS**

		ssh -i kube-key core@15.126.200.248
			
		kubecfg list minions 
		
		Results:

			$ ssh -i ../kube-key core@15.126.200.248
			
			The authenticity of host '15.126.200.248 (15.126.200.248)' can't be established.
			RSA key fingerprint is fe:b1:a0:6f:3b:60:e7:3c:26:30:98:4a:86:24:99:d8.
			Are you sure you want to continue connecting (yes/no)? yes
			Warning: Permanently added '15.126.200.248' (RSA) to the list of known hosts.
			CoreOS (stable)
			core@kube-master ~ $ kubecfg list minions
			Minion identifier
			----------
			10.0.0.40
			10.0.0.41
			10.0.0.42
			
	**Windows**
	
		coming soon
	
8. After verifying all the nodes are there, you are ready to rock and roll. The next step would be to deploy a [sample application](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/examples/guestbook/README.md
) to the cluster. Happy containerizing!

## License ##

Copyright 2014 Hewlett-Packard

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
