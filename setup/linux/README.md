#Install hpcloud-kubesetup on Linux

###Install

To install hpcloud-kubesetup and Kubernetes kubectl on to Linux run the following command:

**From a bash shell:**

    bash <(curl -Ls https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/linux/kubernetes-tools-install.sh)

This will install the following files:
* /usr/local/kubernetes/hpcloud-kubesetup-linux.zip
* /usr/local/kubernetes/linux/hpcloud-kubesetup
* /usr/local/kubernetes/linux/kubectl
* /usr/local/kubernetes/linux/kubesetup.yml
* /usr/local/kubernetes/linux/LICENSE
* /usr/local/kubernetes/linux/README

Working copy of the installation configuration
* ~/kubernetes/kubesetup.yml

###Uninstall

To uninstall hpcloud-kubesetup and Kubernetes kubectl run the following command:

**From a bash shell:**

    bash <(curl -Ls https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/linux/kubernetes-tools-uninstall.sh)
