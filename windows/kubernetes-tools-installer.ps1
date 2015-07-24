# ==============================================================================
# kubernetes-tools-installer.ps1 
#
# installs hpcloud-kubesetup.exe and Kubernetes kubectl.exe
#
# https://github.com/hpcloud/hpcloud-kubesetup/windows/kubernetes-tools-installer.ps1
#
# © Copyright 2015 Hewlett-Packard Development Company, L.P.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use 
# this file except in compliance with the License. You may obtain a copy of the 
# License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed 
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR 
# CONDITIONS OF ANY KIND, either express or implied. See the License for the 
# specific language governing permissions and limitations under the License.
# ==============================================================================
Write-Host "kubernetes-tools-installer.ps1`n"

$installPath = 'c:\kube'
if (Test-Path -Path $installPath) {
	Write-Host "Warning   : InstallPath $installPath already exists, overwriting existing files`n" -foregroundcolor red -backgroundcolor yellow
}
else {
  mkdir $installPath -Force | Out-Null
}

Write-Host "Install to: $installPath`n"

function Download-File {
param (
  [string]$url,
  [string]$file
 )
  Write-Host "Download  : $url"
  Write-Host "to        : $file`n"
  $downloader = new-object System.Net.WebClient
  $downloader.Proxy.Credentials=[System.Net.CredentialCache]::DefaultNetworkCredentials;
  $downloader.DownloadFile($url, $file)
}

Download-File "https://github.com/gertd/hpcloud-kubesetup/raw/master/bin/hpcloud-kubesetup-windows.zip" (Join-Path $InstallPath 'hpcloud-kubesetup-windows.zip')
Download-File "https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/windows/amd64/kubectl.exe" (Join-Path $InstallPath 'kubectl.exe')

Function Unzip-File {
param (
  [string]$sourceFile,
  [string]$targetPath
 )
  Write-Host "Unzip     : $sourceFile"
  Write-Host "to        : $targetPath`n"
  $fileSystemAssemblyPath = Join-Path ([System.Runtime.InteropServices.RuntimeEnvironment]::GetRuntimeDirectory()) 'System.IO.Compression.FileSystem.dll'
  Add-Type -Path $fileSystemAssemblyPath
  [System.IO.Compression.ZipFile]::ExtractToDirectory($sourceFile, $targetPath)
}

Unzip-File (Join-Path $InstallPath 'hpcloud-kubesetup-windows.zip') $InstallPath
Move-Item -Path ((Join-Path $InstallPath 'windows')+'\*') -Destination $InstallPath -Force
Remove-Item -Path (Join-Path $InstallPath 'windows') -Force

Write-Host "Finished"