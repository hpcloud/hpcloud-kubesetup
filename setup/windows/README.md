#Install hpcloud-kubesetup on Windows

###Install

To install hpcloud-kubesetup.exe and Kubernetes kubectl.exe on to Windows run one of the following commands

**From command shell (cmd.exe):**

    @powershell -NoProfile -ExecutionPolicy Bypass -Command "iex ((New-Object Net.WebClient).DownloadString('https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/windows/kubernetes-tools-install.ps1'))"

**From PowerShell:**

    iex ((New-Object Net.WebClient).DownloadString('https://raw.githubusercontent.com/hpcloud/hpcloud-kubesetup/master/setup/windows/kubernetes-tools-install.ps1'))

This will install the following files in the **c:\kube** directory:
* hpcloud-kubesetup-windows.zip
* hpcloud-kubesetup.exe
* kubectl.exe
* kubesetup.yml
* LICENSE
* README
