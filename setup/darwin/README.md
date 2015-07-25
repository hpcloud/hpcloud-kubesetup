#Install hpcloud-kubesetup on MacOSX

###Install

To install hpcloud-kubesetup and Kubernetes kubectl on to MacOSX run the following command:

**From a bash shell:**

    bash <(curl -Ls https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/darwin/kubernetes-tools-install.sh)

This will install the following files:
* /usr/local/kubernetes/hpcloud-kubesetup-darwin.zip
* /usr/local/kubernetes/darwin/hpcloud-kubesetup
* /usr/local/kubernetes/darwin/kubectl
* /usr/local/kubernetes/darwin/kubesetup.yml
* /usr/local/kubernetes/darwin/LICENSE
* /usr/local/kubernetes/darwin/README

Working copy of the installation configuration
* ~/kubernetes/kubesetup.yml

###Uninstall

To uninstall hpcloud-kubesetup and Kubernetes kubectl run the following command:

**From a bash shell:**

    bash <(curl -Ls https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/darwin/kubernetes-tools-uninstall.sh)
