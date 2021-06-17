package k3sair

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
)

func init() {

	k3sInstallCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("arch", "", "", "Enter the target sever os architecture (amd64 supported atm)")
	installCmd.Flags().StringP("base", "", "", "Enter the on site proxy repository url (e.g Artifactory)")
	installCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	installCmd.Flags().IP("ip", nil, "Public IP of node")
	installCmd.Flags().StringP("user", "", "root", "Username for SSH login (Default: root")
	installCmd.Flags().BoolP("sudo", "", true, " Use sudo for installation. (Default: true)")

}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install k3s on a server via SSH",
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
	ip, _ := cmd.Flags().GetIP("ip")
	arch, _ := cmd.Flags().GetString("arch")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")

	fmt.Println(fmt.Sprintf("Downloading %s scripts and binaries", color.BlueString("k3s")))
	air := airgap.NewAirGap(base, arch, key, ip.String(), "", user, sudo)
	err := air.DownloadAirGap()
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Bootstraping %s cluster", color.BlueString("k3s")))
	err = air.Install()
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Downloading %s kubeconfig from %s", color.BlueString("k3s"),
		color.GreenString(ip.String())))
	err = air.GetKubeConfig()
	if err != nil {
		return err
	}
	return err
}
