#!/bin/bash

# execute neutron net-external-list to determin value
EXTERNAL_NETWORK_NAME="Ext-Net"

neutron net-create kube-net
neutron subnet-create --name kube-subnet --allocation-pool start=192.168.1.100,end=192.168.1.200 kube-net 192.168.1.0/24
neutron router-create kube-router
neutron router-gateway-set kube-router ${EXTERNAL_NETWORK_NAME}
neutron router-interface-add kube-router kube-subnet
