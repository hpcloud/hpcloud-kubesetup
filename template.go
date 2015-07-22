package main

import (
	"text/template"
)

/*
  Adopted from https://github.com/kelseyhightower/kubeconfig/template.go
  License: https://github.com/kelseyhightower/kubeconfig/blob/master/LICENSE

  Copyright (c) 2014 Kelsey Hightower
*/
var masterTmpl = template.Must(template.New("node").Parse(`#cloud-config

write_files:
  - path: /opt/bin/waiter.sh
    owner: root
    permissions: 0755
    content: |
      #! /usr/bin/bash
      until curl http://127.0.0.1:2379/v2/machines; do sleep 2; done

coreos:
  etcd2:
    name: master
    initial-cluster-token: k8s_etcd
    initial-cluster: master=http://{{.ip}}:2380
    listen-peer-urls: http://{{.ip}}:2380,http://localhost:2380
    initial-advertise-peer-urls: http://{{.ip}}:2380
    listen-client-urls: http://{{.ip}}:2379,http://localhost:2379
    advertise-client-urls: http://{{.ip}}:2379
  fleet:
    etcd_servers: http://localhost:2379
    metadata: k8srole=master
  flannel:
    etcd_endpoints: http://localhost:2379
  locksmithd:
    endpoint: http://localhost:2379
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
    - name: etcd2-waiter.service
      command: start
      content: |
        [Unit]
        Description=etcd waiter
        Wants=network-online.target
        Wants=etcd2.service
        After=etcd2.service
        After=network-online.target
        Before=flanneld.service fleet.service locksmithd.service

        [Service]
        ExecStart=/usr/bin/bash /opt/bin/waiter.sh
        RemainAfterExit=true
        Type=oneshot
    - name: flanneld.service
      command: start
      drop-ins:
        - name: 50-network-config.conf
          content: |
            [Service]
            ExecStartPre=-/usr/bin/etcdctl mk /coreos.com/network/config '{"Network": "10.244.0.0/16", "Backend": {"Type": "vxlan"}}'
    - name: docker-cache.service
      command: start
      content: |
        [Unit]
        Description=Docker cache proxy
        Requires=early-docker.service
        After=early-docker.service
        Before=early-docker.target

        [Service]
        Restart=always
        TimeoutStartSec=0
        RestartSec=5
        Environment=TMPDIR=/var/tmp/
        Environment=DOCKER_HOST=unix:///var/run/early-docker.sock
        ExecStartPre=-/usr/bin/docker kill docker-registry
        ExecStartPre=-/usr/bin/docker rm docker-registry
        ExecStartPre=/usr/bin/docker pull quay.io/devops/docker-registry:latest
        # GUNICORN_OPTS is an workaround for
        # https://github.com/docker/docker-registry/issues/892
        ExecStart=/usr/bin/docker run --rm --net host --name docker-registry \
            -e STANDALONE=false \
            -e GUNICORN_OPTS=[--preload] \
            -e MIRROR_SOURCE=https://registry-1.docker.io \
            -e MIRROR_SOURCE_INDEX=https://index.docker.io \
            -e MIRROR_TAGS_CACHE_TTL=1800 \
            quay.io/devops/docker-registry:latest
    - name: docker.service
      drop-ins:
        - name: 51-docker-mirror.conf
          content: |
            [Unit]
            # making sure that docker-cache is up and that flanneld finished
            # startup, otherwise containers won't land in flannel's network...
            Requires=docker-cache.service
            After=docker-cache.service

            [Service]
            Environment=DOCKER_OPTS='--registry-mirror=http://{{.ip}}:5000'
    - name: get-kubectl.service
      command: start
      content: |
        [Unit]
        Description=Get kubectl client tool
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=network-online.target
        After=network-online.target

        [Service]
        ExecStart=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl
        ExecStart=/usr/bin/chmod +x /opt/bin/kubectl
        Type=oneshot
        RemainAfterExit=true
    - name: kube-apiserver.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes API Server
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=etcd2-waiter.service
        After=etcd2-waiter.service

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kube-apiserver
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-apiserver
        ExecStart=/opt/bin/kube-apiserver \
        --insecure-bind-address=0.0.0.0 \
        --service-cluster-ip-range=10.100.0.0/16 \
        --etcd-servers=http://localhost:2379
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
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kube-controller-manager
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-controller-manager
        ExecStart=/opt/bin/kube-controller-manager \
        --master=127.0.0.1:8080
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
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kube-scheduler
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-scheduler
        ExecStart=/opt/bin/kube-scheduler \
        --master=127.0.0.1:8080
        Restart=always
        RestartSec=10
    - name: kube-register.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Registration Service
        Documentation=https://github.com/kelseyhightower/kube-register
        Requires=kube-apiserver.service fleet.service
        After=kube-apiserver.service fleet.service

        [Service]
        ExecStartPre=-/usr/bin/wget -nc -O /opt/bin/kube-register https://github.com/kelseyhightower/kube-register/releases/download/v0.0.3/kube-register-0.0.3-linux-amd64
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-register
        ExecStart=/opt/bin/kube-register \
        --metadata=k8srole=node \
        --fleet-endpoint=unix:///var/run/fleet.sock \
        --api-endpoint=http://127.0.0.1:8080
        Restart=always
        RestartSec=10
  update:
    group: alpha
    reboot-strategy: off

ssh_authorized_keys:
    - {{.sshkey}}`))

var nodeTmpl = template.Must(template.New("node").Parse(`#cloud-config

write_files:
  - path: /opt/bin/wupiao
    owner: root
    permissions: 0755
    content: |
      #!/bin/bash
      # [w]ait [u]ntil [p]ort [i]s [a]ctually [o]pen
      [ -n "$1" ] && [ -n "$2" ] && while ! curl --output /dev/null \
        --silent --head --fail \
        http://${1}:${2}; do sleep 1 && echo -n .; done;
      exit $?

coreos:
  etcd2:
    listen-client-urls: http://localhost:2379
    advertise-client-urls: http://0.0.0.0:2379
    initial-cluster: master=http://{{.master}}:2380
    proxy: on
  fleet:
    etcd_servers: http://localhost:2379
    metadata: k8srole=node
  flannel:
    etcd_endpoints: http://localhost:2379
  locksmithd:
    endpoint: http://localhost:2379
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
    - name: flanneld.service
      command: start
    - name: docker.service
      command: start
      drop-ins:
        - name: 50-docker-mirror.conf
          content: |
            [Service]
            Environment=DOCKER_OPTS='--registry-mirror=http://{{.master}}:5000'
    - name: kubelet.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Kubelet
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=network-online.target
        After=network-online.target

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubelet
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kubelet
        # wait for kubernetes master to be up and ready
        ExecStartPre=/opt/bin/wupiao {{.master}} 8080
        ExecStart=/opt/bin/kubelet \
        --api-servers={{.master}}:8080 \
        --hostname-override={{.ip}}
        Restart=always
        RestartSec=10
    - name: kube-proxy.service
      command: start
      content: |
        [Unit]
        Description=Kubernetes Proxy
        Documentation=https://github.com/GoogleCloudPlatform/kubernetes
        Requires=network-online.target
        After=network-online.target

        [Service]
        ExecStartPre=/usr/bin/wget -N -P /opt/bin https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kube-proxy
        ExecStartPre=/usr/bin/chmod +x /opt/bin/kube-proxy
        # wait for kubernetes master to be up and ready
        ExecStartPre=/opt/bin/wupiao {{.master}} 8080
        ExecStart=/opt/bin/kube-proxy \
        --master=http://{{.master}}:8080
        Restart=always
        RestartSec=10

  update:
    group: alpha
    reboot-strategy: off

ssh_authorized_keys:
    - {{.sshkey}}`))
