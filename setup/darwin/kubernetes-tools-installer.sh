#!/bin/bash

mkdir -p ~/kube

wget https://github.com/hpcloud/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-darwin.zip \
-O ~/kube/hpcloud-kubesetup-darwin.zip

unzip hpcloud-kubesetup-darwin.zip -d ~/kube/

sudo mv ~/kube/darwin/hpcloud-kubesetup /usr/local/bin/hpcloud-kubesetup
sudo chmod +x /usr/local/bin/hpcloud-kubesetup

wget https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/darwin/amd64/kubectl \
-O ~/kube/kubectl

sudo mv ~/kube/kubectl /usr/local/bin/kubectl
sudo chmod +x /usr/local/bin/kubectl

