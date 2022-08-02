package common

import (
	"io/ioutil"
	"strings"
)

const (
	Cmd1 = `chmod +x /tmp/install.sh
sudo mkdir -p /var/lib/rancher/k3s/agent/images/
sudo mkdir -p /etc/rancher/k3s/
`
	Cmd2 = `chmod +x %s 
sudo cp %s /usr/local/bin/k3s
`
	Cmd3                          = "sudo cp %s /var/lib/rancher/k3s/agent/images/"
	Cmd4                          = "sudo cat /var/lib/rancher/k3s/server/node-token"
	JoinCmd                       = `INSTALL_K3S_SKIP_DOWNLOAD=true K3S_TOKEN=%s K3S_URL=https://%s:%d INSTALL_K3S_EXEC="%s" /tmp/install.sh`
	InstallCmd                    = `INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_EXEC="%s" /tmp/install.sh`
	KubeConfigCmd                 = "sudo cat /etc/rancher/k3s/k3s.yaml"
	Amd64BinaryName               = "k3s-airgap-images-amd64.tar.gz"
	ArmBinaryName                 = "k3s-airgap-images-arm.tar.gz"
	K3sBinary                     = "k3s"
	K3sArmBinary                  = "k3s-arm64"
	K3sYaml                       = "k3s.yaml"
	InstallRegistriesYamlLocation = "sudo cp /tmp/registries.yaml /etc/rancher/k3s/registries.yaml"
	TmpInstallScript              = "%s/install.sh"
	TmpRegistriesYaml             = "%s/registries.yaml"
	ARM64                         = "aarch64"
)

func WriteFile(filename, content string) error {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func CheckSudo(sudo bool, cmd string) string {
	var command = cmd
	if !sudo {
		command = strings.Replace(command, "sudo", "", -1)
	}
	return command
}
