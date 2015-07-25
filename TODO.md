TODO list hpcloud-kubesetup
===========================

1.  Validate provided availabilityZone before create server call
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
13. More input validation flavor name, network name, network ip in range of network name, network name does not have to be unique, allow for id
14. Rename install->create uninstall->delete, to align with add & remove
15. Improve/cleanup debug output feed
