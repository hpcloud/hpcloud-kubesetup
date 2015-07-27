TODO list hpcloud-kubesetup
===========================

1.  ~~Validate provided availabilityZone before create server call~~
2.  Add node to cluster
3.  Remove node from cluster
4.  Create security group for external communication kubernetes-external
5.  Create security group for internal communication kubernetes-internal
6.  Enable status command line option for displaying cluster status at IaaS level
7.  Add --debug to file
8.  Use DHCP assigned network addresses for Nodes
9.  Determine master IP address based on network and first available IP in range
10. Install kubectl on client which is running kubesetup
11. Create the client ssh tunel to the master (ssh -f -nNT -L 8080:127.0.0.1:8080 core@<master-public-ip>)
12. Set http proxy information on nodes using CloudInit
13. More input validation ~~flavor name, network name~~, network ip in range of network name, ~~network name does not have to be unique~~, allow for network id
14. Rename install->create uninstall->delete, to align with add & remove
15. Improve/cleanup debug output feed
16. Assign cluster id to master node, add cluster id to all nodes in nova
17. Allow for id input besides names for all inputs
18. Rework command line arguments
    * create - creates the cluster, aka the master node
    * add - adds nodes to the cluster providing a name, optional IP, default to auto assigned DHCP IP for node
    * remove - remove nodes from the cluster by name or IP
    * delete - deletes the cluster, and any remaining nodes when using -force
    * list - list the members of the cluster
    * status - report status of each node in cluster
19. Rework config file
  * Objective is to have a single config file for all operation against a cluster
  * The config file will have 3 sections: cluster configuration, template sections for master & node

cluster:
  name: k8-cluster
  sshkey: kube-key
  network: kube-net
  master-ip: 192.168.1.140

templates:
  master:
    image-name: CoreOS
    image-id:
    flavor-name: standard.medium
    flavor-id:
    availabilityZone: AZ2
  node:
    image-name: CoreOS
    image-id:
    flavor-name: standard.medium
    flavor-id:
    availabilityZone: AZ2

20. Only generate two CloudInit files, one for master and one for node, instead of one per machine
21. Add command switch to bypass creating cloudinit files, enabling manual changes to the cloudinit files
