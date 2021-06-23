package k3sair

import (
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
)

func init() {

	k3sInstallCmd.AddCommand(installCmd)

	installCmd.Flags().String("arch", "", "Enter the target sever os architecture (amd64 supported atm)")
	installCmd.Flags().String("base", "", "Enter the on site proxy repository url (e.g Artifactory)")
	installCmd.Flags().String("ssh-key", "", "The ssh key to use for remote login")
	installCmd.Flags().String("ip", "", "Public ip or FQDN of node")
	installCmd.Flags().String("user", "root", "Username for SSH login (Default: root")
	installCmd.Flags().Bool("sudo", true, "Use sudo for installation. (Default: true)")
	installCmd.Flags().String("mirror", "", "Mirrored Registry. (Default: '')")
	installCmd.Flags().String("tls-san", "", "Add additional hostname or IP as a Subject Alternative Name in the TLS cert")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "InstallControlPlaneNode k3s on a server via SSH",
	Example: `k3sair install  \
    --arch amd64
    --base https//artifactory.local/generic/
    --ssh-key ~/.ssh/id_rsa
    --ip 127.0.0.1
`,
	RunE:          runInstall,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runInstall(cmd *cobra.Command, _ []string) error {

	base, _ := cmd.Flags().GetString("base")
	key, _ := cmd.Flags().GetString("ssh-key")
	ip, _ := cmd.Flags().GetString("ip")
	arch, _ := cmd.Flags().GetString("arch")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")
	mirror, _ := cmd.Flags().GetString("mirror")
	tlsSAN, _ := cmd.Flags().GetString("tls-san")

	air := airgap.NewAirGap(base, arch, key, ip, user, sudo)
	err := air.InstallAirGapFiles(mirror)
	if err != nil {
		return err
	}

	if len(tlsSAN) > 0 {
		air.AddServerOptions("--tls-san " + tlsSAN)
	}

	err = air.InstallControlPlaneNode()
	if err != nil {
		return err
	}

	err = air.GetKubeConfig()
	if err != nil {
		return err
	}
	return err
}
