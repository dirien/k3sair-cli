package k3sair

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Version as per git repo
	Version string

	// GitCommit as per git repo
	GitCommit string
)

func init() {
	k3sInstallCmd.AddCommand(versionCmd)
}

var k3sInstallCmd = &cobra.Command{
	Use:   "k3sair",
	Short: "Air-Gap InstallControlPlaneNode of a k3s cluster",
	Run:   runK3sair,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

func Execute(version, gitCommit string) error {

	Version = version
	GitCommit = gitCommit

	if err := k3sInstallCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runK3sair(cmd *cobra.Command, _ []string) {
	printLogo()
	cmd.Help()
}

func printLogo() {
	logo := color.GreenString(figletStr)
	fmt.Println(logo)
}

const figletStr = `
██╗  ██╗██████╗ ███████╗ █████╗ ██╗██████╗ 
██║ ██╔╝╚════██╗██╔════╝██╔══██╗██║██╔══██╗
█████╔╝  █████╔╝███████╗███████║██║██████╔╝
██╔═██╗  ╚═══██╗╚════██║██╔══██║██║██╔══██╗
██║  ██╗██████╔╝███████║██║  ██║██║██║  ██║
╚═╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝╚═╝  ╚═╝
`
