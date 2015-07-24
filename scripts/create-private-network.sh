#!/bin/bash

EXTERNAL_NETWORK_NAME="external"

neutron net-create als-network
neutron subnet-create --name als-subnet --allocation-pool start=192.168.15.10,end=192.168.15.250 --dns-nameserver 206.164.176.23 als-network 192.168.15.0/24
neutron router-create als-router
neutron router-gateway-set als-router ${EXTERNAL_NETWORK_NAME}
neutron router-interface-add als-router als-subnet
