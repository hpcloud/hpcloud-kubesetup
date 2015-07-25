#!/bin/bash

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
