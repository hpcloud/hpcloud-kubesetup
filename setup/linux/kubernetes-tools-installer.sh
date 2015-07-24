#!/bin/bash

mkdir -p ~/kube

wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-linux.zip \
-O ~/kube/hpcloud-kubesetup-linux.zip

unzip hpcloud-kubesetup-linux.zip -d ~/kube/

sudo mv ~/kube/linux/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup
sudo chmod +x /usr/local/bin/hpcloud-kubesetup

wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl \
-O ~/kube/kubectl

sudo mv ~/kube/kubectl /usr/local/bin/kubectl
sudo chmod +x /usr/local/bin/kubectl
