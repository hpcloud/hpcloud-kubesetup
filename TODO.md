TODO list hpcloud-kubesetup
===========================

1. Validate provided availabilityZone before create server call
2. Create security group for external communication kubernetes-external
3. Create security group for internal communication kubernetes-internal
4. Enable status command line option for displaying cluster status at IaaS level
5. Send --debug output to separate output channel
6. Use DHCP assigned network addresses for Nodes
7. Determine master IP address based on network and first available IP in range
8. Install kubectl on client which is running kubesetup
9. Create the client ssh tunel to the master (ssh -f -nNT -L 8080:127.0.0.1:8080 core@<master-public-ip>)
