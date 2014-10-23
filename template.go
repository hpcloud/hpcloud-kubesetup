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

hostname: {{.hostname}}
coreos:
  etcd:
    discovery: {{.discovery}}
    name: {{.hostname}}
    addr: {{.ip}}:4001
    bind-addr: 0.0.0.0
    peer-addr: {{.ip}}:7001
{{if ne .role "master"}}
    peers: {{.peers}}
{{end}}
    peer-heartbeat-interval: 250
    peer-election-timeout: 1000
  units:
    - name: static.network
      command: start
      content: |
        [Match]
        Name=ens33

        [Network]
        Address={{.ip}}/24
        DNS={{.dns}}
        Gateway={{.gateway}}
    - name: cbr0.netdev
      command: start
      content: |
        [NetDev]
        Kind=bridge
        Name=cbr0
    - name: cbr0.network
      command: start
      content: |
        [Match]
        Name=cbr0

        [Network]
        Address={{.subnet}}

        [Route]
        Destination=10.0.0.0/8
        Gateway=0.0.0.0
    - name: cbr0-interface.network
      command: start
      content: |
        [Match]
        Name=ens34

        [Network]
        Bridge=cbr0
    - name: nat.service
      command: start
      content: |
        [Unit]
        Description=NAT non container traffic

        [Service]
        ExecStart=/usr/sbin/iptables -t nat -A POSTROUTING -o ens33 -j MASQUERADE ! -d 10.0.0.0/8
        RemainAfterExit=yes
        Type=oneshot
    - name: etcd.service
      command: start
    - name: fleet.service
      command: start
    - name: docker.service
      command: start
      content: |
        [Unit]
        After=network.target
        Description=Docker Application Container Engine
        Documentation=http://docs.docker.io

        [Service]
        ExecStartPre=/bin/mount --make-rprivate /
        ExecStart=/usr/bin/docker -d -s=btrfs -H fd:// -b cbr0 --iptables=false

        [Install]
        WantedBy=multi-user.target
    - name: download-kubernetes.service
      command: start
      content: |
        [Unit]
        After=network-online.target
{{if eq .role "master"}}
        Before=apiserver.service
        Before=controller-manager.service
{{end}}
        Before=kubelet.service
        Before=proxy.service
        Description=Download Kubernetes Binaries
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=network-online.target

        [Service]
{{if eq .role "master"}}
        ExecStart=/usr/bin/wget -N -P /opt/bin http://storage.googleapis.com/kubernetes/apiserver
        ExecStart=/usr/bin/wget -N -P /opt/bin http://storage.googleapis.com/kubernetes/controller-manager
        ExecStart=/usr/bin/wget -N -P /opt/bin http://storage.googleapis.com/kubernetes/kubecfg
{{end}}
        ExecStart=/usr/bin/wget -N -P /opt/bin http://storage.googleapis.com/kubernetes/kubelet
        ExecStart=/usr/bin/wget -N -P /opt/bin http://storage.googleapis.com/kubernetes/proxy
{{if eq .role "master"}}
        ExecStart=/usr/bin/chmod +x /opt/bin/apiserver
        ExecStart=/usr/bin/chmod +x /opt/bin/controller-manager
        ExecStart=/usr/bin/chmod +x /opt/bin/kubecfg
{{end}}
        ExecStart=/usr/bin/chmod +x /opt/bin/kubelet
        ExecStart=/usr/bin/chmod +x /opt/bin/proxy
        RemainAfterExit=yes
        Type=oneshot
{{if eq .role "master"}}
    - name: apiserver.service
      command: start
      content: |
        [Unit]
        ConditionFileIsExecutable=/opt/bin/apiserver
        Description=Kubernetes API Server
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes

        [Service]
        ExecStart=/opt/bin/apiserver \
        --address=127.0.0.1 \
        --port=8080 \
        --etcd_servers=http://127.0.0.1:4001 \
        --machines={{.machines}} \
        --logtostderr=true
        Restart=always
        RestartSec=2

        [Install]
        WantedBy=multi-user.target
    - name: controller-manager.service
      command: start
      content: |
        [Unit]
        ConditionFileIsExecutable=/opt/bin/controller-manager
        Description=Kubernetes Controller Manager
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes

        [Service]
        ExecStart=/opt/bin/controller-manager \
        --etcd_servers=http://127.0.0.1:4001 \
        --master=127.0.0.1:8080 \
        --logtostderr=true
        Restart=always
        RestartSec=2

        [Install]
        WantedBy=multi-user.target
{{end}}
    - name: kubelet.service
      command: start
      content: |
        [Unit]
        ConditionFileIsExecutable=/opt/bin/kubelet
        Description=Kubernetes Kubelet
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes

        [Service]
        ExecStart=/opt/bin/kubelet \
        --address=0.0.0.0 \
        --port=10250 \
        --hostname_override={{.ip}} \
        --etcd_servers=http://127.0.0.1:4001 \
        --logtostderr=true
        Restart=always
        RestartSec=2

        [Install]
        WantedBy=multi-user.target
    - name: proxy.service
      command: start
      content: |
        [Unit]
        ConditionFileIsExecutable=/opt/bin/proxy
        Description=Kubernetes Proxy
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes

        [Service]
        ExecStart=/opt/bin/proxy --etcd_servers=http://127.0.0.1:4001 --logtostderr=true
        Restart=always
        RestartSec=2

        [Install]
        WantedBy=multi-user.target
  update:
    group: alpha
    reboot-strategy: off
ssh_authorized_keys:
    - {{.sshkey}}`))
