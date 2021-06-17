package k3sair

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/airgap"
	"github.com/spf13/cobra"
	"log"
)

func init() {

	k3sInstallCmd.AddCommand(joinCmd)

	joinCmd.Flags().StringP("arch", "", "", "Enter the target sever os architecture (amd64 supported atm)")
	joinCmd.Flags().StringP("base", "", "", "Enter the on site proxy repository url (e.g Artifactory)")
	joinCmd.Flags().StringP("ssh-key", "", "", "The ssh key to use for remote login")
	joinCmd.Flags().IP("ip", nil, "Public IP of node")
	joinCmd.Flags().StringP("user", "", "root", "Username for SSH login (Default: root")
	joinCmd.Flags().BoolP("sudo", "", true, " Use sudo for installation. (Default: true)")
	joinCmd.Flags().IP("control-plane-ip", nil, "Public IP of an existing k3s server")

}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Install the k3s agent on a remote host and join it to an existing server",
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
	ip, _ := cmd.Flags().GetIP("ip")
	controlPlaneIp, _ := cmd.Flags().GetIP("control-plane-ip")
	user, _ := cmd.Flags().GetString("user")
	sudo, _ := cmd.Flags().GetBool("sudo")

	fmt.Println(fmt.Sprintf("Downloading %s scripts and binaries\n", color.BlueString("k3s")))
	air := airgap.NewAirGap(base, arch, key, ip.String(), controlPlaneIp.String(), user, sudo)
	err := air.DownloadAirGap()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Joining existing %s cluster %s\n", color.BlueString("k3s"), color.GreenString(controlPlaneIp.String())))
	err = air.Join()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
