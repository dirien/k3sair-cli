package k3sair

import (
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
	"log"
)

func init() {

	k3sInstallCmd.AddCommand(joinCmd)

	joinCmd.Flags().StringP("arch", "", "", "Enter the target sever os architecture (amd64 supported atm)")
	joinCmd.Flags().StringP("base", "", "", "Enter the on site proxy repository url (e.g Artifactory)")
	joinCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	joinCmd.Flags().StringP("ip", "", "", "Public ip or FQDN of node")
	joinCmd.Flags().StringP("user", "", "root", "Username for SSH login (Default: root")
	joinCmd.Flags().BoolP("sudo", "", true, " Use sudo for installation. (Default: true)")
	joinCmd.Flags().StringP("control-plane-ip", "", "", "Public ip or FQDN of an existing k3s server")
	joinCmd.Flags().StringP("mirror", "", "", "Mirrored Registry. (Default: '')")

}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "InstallControlPlaneNode the k3s agent on a remote host and join it to an existing server",
	Example: `k3sair join  \
    --arch amd64
    --base https//artifactory.local/generic/
	--ssh-key ~/.ssh/id_rsa
    --ip 127.0.0.1
	--control-plane-ip 127.0.0.2`,
	RunE:          joinCreate,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func joinCreate(cmd *cobra.Command, _ []string) error {
	base, _ := cmd.Flags().GetString("base")
	key, _ := cmd.Flags().GetString("ssh-key")
	arch, _ := cmd.Flags().GetString("arch")
	ip, _ := cmd.Flags().GetString("ip")
	controlPlaneIp, _ := cmd.Flags().GetString("control-plane-ip")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")
	mirror, _ := cmd.Flags().GetString("mirror")

	air := airgap.NewAirGap(base, arch, key, ip, user, sudo)
	err := air.InstallAirGapFiles(mirror)
	if err != nil {
		log.Fatal(err)
	}

	controlPlane := airgap.NewAirGap(base, arch, key, controlPlaneIp, user, sudo)
	token, err := controlPlane.GetNodeToken()
	if err != nil {
		return err
	}

	err = air.InstallWorkerNode(controlPlaneIp, token)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
