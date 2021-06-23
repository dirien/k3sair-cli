package k3sair

import (
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
)

func init() {

	k3sInstallCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("arch", "", "", "Enter the target sever os architecture (amd64 supported atm)")
	installCmd.Flags().StringP("base", "", "", "Enter the on site proxy repository url (e.g Artifactory)")
	installCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	installCmd.Flags().StringP("ip", "", "", "Public ip or FQDN of node")
	installCmd.Flags().StringP("user", "", "root", "Username for SSH login (Default: root")
	installCmd.Flags().BoolP("sudo", "", true, " Use sudo for installation. (Default: true)")
	installCmd.Flags().StringP("mirror", "", "", "Mirrored Registry. (Default: '')")

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

	air := airgap.NewAirGap(base, arch, key, ip, user, sudo)
	err := air.InstallAirGapFiles(mirror)

	if err != nil {
		return err
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
