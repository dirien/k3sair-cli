package k3sair

import (
	"errors"
	"log"

	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
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
	joinCmd.Flags().Uint("control-plane-port", 22, "The ssh port to use")
	joinCmd.Flags().Uint("k3s-api-port", 6443, "The kube api server port.")
	joinCmd.Flags().StringP("mirror", "", "", "Mirrored Registry. (Default: '')")
	joinCmd.Flags().Uint("port", 22, "The ssh port to use")
	joinCmd.Flags().String("additional-k3s-exec-flags", "", "Add additional k3s exec flags, separate with space")
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
	controlPlaneIP, _ := cmd.Flags().GetString("control-plane-ip")
	controlPlanePort, _ := cmd.Flags().GetUint("control-plane-port")
	k3sAPIPort, _ := cmd.Flags().GetUint("k3s-api-port")
	port, _ := cmd.Flags().GetUint("port")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")
	mirror, _ := cmd.Flags().GetString("mirror")
	additionalK3sExecFlags, _ := cmd.Flags().GetString("additional-k3s-exec-flags")

	if len(base) == 0 {
		return errors.New("on-site proxy repository must be provided")
	}

	air := airgap.NewAirGap(base, arch, key, ip, user, port, sudo)
	err := air.InstallAirGapFiles(mirror)
	if err != nil {
		log.Fatal(err)
	}

	controlPlane := airgap.NewAirGap(base, arch, key, controlPlaneIP, user, controlPlanePort, sudo)
	token, err := controlPlane.GetNodeToken()
	if err != nil {
		return err
	}

	if len(additionalK3sExecFlags) > 0 {
		air.AddServerOptions(additionalK3sExecFlags)
	}

	err = air.InstallWorkerNode(controlPlaneIP, token, k3sAPIPort)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
