package k3sair

import (
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
)

func init() {
	k3sInstallCmd.AddCommand(kubeConfigCmd)

	kubeConfigCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	kubeConfigCmd.Flags().StringP("ip", "", "", "Public ip or FQDN of node")
	kubeConfigCmd.Flags().Uint("port", 22, "The ssh port to use")
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
	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetUint("port")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")

	air := airgap.NewAirGap("", "", key, ip, user, port, sudo)
	err := air.GetKubeConfig()
	if err != nil {
		return err
	}
	return err
}
