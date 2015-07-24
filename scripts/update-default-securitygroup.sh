#!/bin/bash

neutron security-group-rule-create --protocol tcp --direction ingress --port-range-min 22 --port-range-max 22 default
neutron security-group-rule-create --protocol tcp --direction ingress --port-range-min 80 --port-range-max 80 default
neutron security-group-rule-create --protocol tcp --direction ingress --port-range-min 8080 --port-range-max 8080 default
