package common

import "strings"

const (
	Cmd1 = `chmod +x install.sh
sudo mkdir -p /var/lib/rancher/k3s/agent/images/
`
	Cmd2 = `chmod +x %s 
sudo cp %s /opt/bin
`
	Cmd3                  = "sudo cp %s /var/lib/rancher/k3s/agent/images/"
	Cmd4                  = "sudo cat /var/lib/rancher/k3s/server/node-token"
	JoinCmd               = "INSTALL_K3S_SKIP_DOWNLOAD=true K3S_TOKEN=%s "
	JoinCmdPart2          = "%s K3S_URL=https://%s:6443 ./install.sh"
	InstallCmd            = "INSTALL_K3S_SKIP_DOWNLOAD=true ./install.sh"
	KubeConfigCmd         = "sudo cat /etc/rancher/k3s/k3s.yaml\n"
	Amd64BinaryName       = "k3s-airgap-images-amd64.tar.gz"
	ArmBinaryName         = "k3s-airgap-images-arm.tar.gz"
	K3sBinary             = "k3s"
	K3sYaml               = "k3s.yaml"
	InstallScriptLocation = "/home/%s/install.sh"
)

func CheckSudo(sudo bool, cmd string) string {
	var command = cmd
	if !sudo {
		command = strings.Replace(command, "sudo", "", -1)
	}
	return command
}
