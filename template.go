package main

import (
  "text/template"
)

/*
  Adopted from https://github.com/kelseyhightower/kubeconfig/template.go
  License: https://github.com/kelseyhightower/kubeconfig/blob/master/LICENSE

  Copyright (c) 2014 Kelsey Hightower
*/

var nodeTmpl = template.Must(template.New("node").Parse(`#cloud-config
---
{{if eq .role "master"}}
write_files:
- path: /opt/bin/waiter.sh
  owner: root
  content: |
    #! /usr/bin/bash
    until curl http://127.0.0.1:4001/v2/machines; do sleep 2; done
coreos:
  units:
    - name: setup-network-environment.service
      command: start
      content: |
        [Unit]
        Description=Setup Network Environment
        Documentation=https://github.com/kelseyhightower/setup-network-environment
        Requires=network-online.target
        After=network-online.target

        [Service]
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/k8s/setup-network-environment
        ExecStartPre=/usr/bin/chmod +x /opt/bin/setup-network-environment
        ExecStart=/opt/bin/setup-network-environment
        RemainAfterExit=yes
        Type=oneshot
    - name: etcd.service
      command: start
      content: |
        [Unit]
        Description=etcd
        Requires=setup-network-environment.service
        After=setup-network-environment.service

        [Service]
        EnvironmentFile=/etc/network-environment
        User=etcd
        PermissionsStartOnly=true
        ExecStart=/usr/bin/etcd \
        --name ${DEFAULT_IPV4} \
        --addr ${DEFAULT_IPV4}:4001 \
        --bind-addr 0.0.0.0 \
        --cluster-active-size 1 \
        --data-dir /var/lib/etcd \
        --http-read-timeout 86400 \
        --peer-addr ${DEFAULT_IPV4}:7001 \
        --snapshot true
        Restart=always
        RestartSec=10s
    - name: fleet.socket
      command: start
      content: |
        [Socket]
        ListenStream=/var/run/fleet.sock
    - name: fleet.service
      command: start
      content: |
        [Unit]
        Description=fleet daemon
        Wants=etcd.service
        After=etcd.service
        Wants=fleet.socket
        After=fleet.socket

        [Service]
        Environment="FLEET_ETCD_SERVERS=http://127.0.0.1:4001"
        Environment="FLEET_METADATA=role=master"
        ExecStart=/usr/bin/fleetd
        Restart=always
        RestartSec=10s
    - name: etcd-waiter.service
      command: start
      content: |
        [Unit]
        Description=etcd waiter
        Wants=network-online.target
        Wants=etcd.service
        After=etcd.service
        After=network-online.target
        Before=flannel.service
        Before=setup-network-environment.service

        [Service]
        ExecStartPre=/usr/bin/chmod +x /opt/bin/waiter.sh
        ExecStart=/usr/bin/bash /opt/bin/waiter.sh
        RemainAfterExit=true
        Type=oneshot
    - name: flannel.service
      command: start
      content: |
        [Unit]
        Wants=etcd-waiter.service
        After=etcd-waiter.service
        Requires=etcd.service
        After=etcd.service
        After=network-online.target
        Wants=network-online.target
        Description=flannel is an etcd backed overlay network for containers

        [Service]
        Type=notify
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/k8s/flanneld
        ExecStartPre=/usr/bin/chmod +x /opt/bin/flanneld
        ExecStartPre=/usr/bin/etcdctl mk /coreos.com/network/config '{"Network":"10.244.0.0/16", "Backend": {"Type": "vxlan"}}'
        ExecStart=/opt/bin/flanneld
    - name: kube-apiserver.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes API Server
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=etcd.service
        After=etcd.service

        [Service]
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kube-apiserver
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-apiserver
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kubecfg
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kubecfg
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kubectl
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kubectl
        ExecStart=/opt/bin/kube-apiserver \
        --address=0.0.0.0 \
        --port=8080 \
        --portal_net=10.100.0.0/16 \
        --etcd_servers=http://127.0.0.1:4001 \
        --public_address_override=$private_ipv4 \
        --logtostderr=true
        Restart=always
        RestartSec=10
    - name: kube-controller-manager.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Controller Manager
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=kube-apiserver.service
        After=kube-apiserver.service

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kube-controller-manager
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-controller-manager
        ExecStart=/opt/bin/kube-controller-manager \
        --master=127.0.0.1:8080 \
        --logtostderr=true
        Restart=always
        RestartSec=10
    - name: kube-scheduler.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Scheduler
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=kube-apiserver.service
        After=kube-apiserver.service

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kube-scheduler
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-scheduler
        ExecStart=/opt/bin/kube-scheduler --master=127.0.0.1:8080
        Restart=always
        RestartSec=10
    - name: kube-register.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Registration Service
        Documentation=https://github.com/kelseyhightower/kube-register
        Requires=kube-apiserver.service
        After=kube-apiserver.service
        Requires=fleet.service
        After=fleet.service

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/k8s/kube-register
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-register
        ExecStart=/opt/bin/kube-register \
        --metadata=role=node \
        --fleet-endpoint=unix:///var/run/fleet.sock \
        --api-endpoint=http://127.0.0.1:8080
        Restart=always
        RestartSec=10
  update:
    group: alpha
    reboot-strategy: off
ssh_authorized_keys:
   - {{.sshkey}}))
{{end}}
{{if eq .role "minion"}}
coreos:
  units:
    - name: etcd.service
      mask: true
    - name: fleet.service
      command: start
      content: |
        [Unit]
        Description=fleet daemon
        Wants=fleet.socket
        After=fleet.socket

        [Service]
        Environment="FLEET_ETCD_SERVERS=http://{{.masterip}}:4001"
        Environment="FLEET_METADATA=role=node"
        ExecStart=/usr/bin/fleetd
        Restart=always
        RestartSec=10s
    - name: flannel.service
      command: start
      content: |
        [Unit]
        After=network-online.target
        Wants=network-online.target
        Description=flannel is an etcd backed overlay network for containers

        [Service]
        Type=notify
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/k8s/flanneld
        ExecStartPre=/usr/bin/chmod +x /opt/bin/flanneld
        ExecStart=/opt/bin/flanneld -etcd-endpoints http://{{.masterip}}:4001
    - name: docker.service
      command: start
      content: |
        [Unit]
        After=flannel.service
        Wants=flannel.service
        Description=Docker Application Container Engine
        Documentation=http://docs.docker.io

        [Service]
        EnvironmentFile=/run/flannel/subnet.env
        ExecStartPre=/bin/mount --make-rprivate /
        ExecStart=/usr/bin/docker -d --bip=${FLANNEL_SUBNET} --mtu=${FLANNEL_MTU} -s=overlay -H fd://

        [Install]
        WantedBy=multi-user.target
    - name: setup-network-environment.service
      command: start
      content: |
        [Unit]
        Description=Setup Network Environment
        Documentation=https://github.com/kelseyhightower/setup-network-environment
        Requires=network-online.target
        After=network-online.target

        [Service]
        ExecStartPre=-/usr/bin/mkdir -p /opt/bin
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/k8s/setup-network-environment
        ExecStartPre=/usr/bin/chmod +x /opt/bin/setup-network-environment
        ExecStart=/opt/bin/setup-network-environment
        RemainAfterExit=yes
        Type=oneshot
    - name: kube-proxy.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Proxy
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=setup-network-environment.service
        After=setup-network-environment.service

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kube-proxy
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-proxy
        ExecStart=/opt/bin/kube-proxy \
        --etcd_servers=http://{{.masterip}}:4001 \
        --logtostderr=true
        Restart=always
        RestartSec=10
    - name: kube-kubelet.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Kubelet
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=setup-network-environment.service
        After=setup-network-environment.service

        [Service]
        EnvironmentFile=/etc/network-environment
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v0.9.1/bin/linux/amd64/kubelet
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kubelet
        ExecStart=/opt/bin/kubelet \
        --address=0.0.0.0 \
        --port=10250 \
        --hostname_override={{.ip}} \
        --etcd_servers=http://{{.masterip}}:4001 \
        --logtostderr=true
        Restart=always
        RestartSec=10
  update:
    group: alpha
    reboot-strategy: off
  ssh_authorized_keys:
    - {{.sshkey}}))
{{end}}
`))
