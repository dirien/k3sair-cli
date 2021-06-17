package k3sair

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
)

func init() {

	k3sInstallCmd.AddCommand(kubeConfigCmd)

	kubeConfigCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	kubeConfigCmd.Flags().IP("ip", nil, "Public IP of node")
	kubeConfigCmd.Flags().StringP("user", "", "root", "Username for SSH login (Default: root")
	kubeConfigCmd.Flags().BoolP("sudo", "", true, " Use sudo for installation. (Default: true)")

}

var kubeConfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get the kubeconfig from the k3s control plane server",
	Example: `k3sair kubeconfig  \
    --ssh-key ~/.ssh/id_rsa
    --ip 127.0.0.1
`,
	RunE:          runKubeConfig,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runKubeConfig(cmd *cobra.Command, _ []string) error {

	key, _ := cmd.Flags().GetString("ssh-key")
	ip, _ := cmd.Flags().GetIP("ip")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")

	fmt.Println(fmt.Sprintf("Downloading %s kubeconfig from %s", color.BlueString("k3s"),
		color.GreenString(ip.String())))
	air := airgap.NewAirGap("", "", key, ip.String(), "", user, sudo)
	err := air.GetKubeConfig()
	if err != nil {
		return err
	}
	return err
}
