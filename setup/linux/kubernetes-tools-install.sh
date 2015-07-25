#!/bin/bash

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
