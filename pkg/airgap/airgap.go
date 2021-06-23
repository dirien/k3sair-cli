package airgap

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/common"
	"github.com/k3sair/pkg/downloader"
	"github.com/k3sair/pkg/embedded"
	"github.com/k3sair/pkg/server"
	"github.com/k3sair/pkg/term"
	"strings"
)

type AirGap struct {
	base                  string
	binary                string
	images                string
	installK3sExec        string
	remoteServer          *server.RemoteServer
	airGapeFileDownloader *downloader.AirGapeFileDownloader
	embeddedFileLoader    *embedded.EmbeddedFileLoader
	color                 *term.Color
	helper                *common.Helper
}

//go:embed install.sh
var installScript string

//go:embed registries.yaml
var registries string

type AirGapped interface {
	InstallAirGapFiles(mirror string) error
	InstallControlPlaneNode() error
	InstallWorkerNode(controlPlaneIp, token string) error
	GetKubeConfig() error
	GetNodeToken() (string, error)
	AddServerOptions(string) *AirGap
}

func (a *AirGap) AddServerOptions(options string) *AirGap {
	a.installK3sExec = fmt.Sprintf("%s %s", a.installK3sExec, options)
	return a
}

func (a *AirGap) GetNodeToken() (string, error) {
	token, err := a.remoteServer.ExecuteCommand(common.Cmd4)
	if err != nil {
		return "", err
	}

	token = strings.TrimSuffix(token, "\n")
	fmt.Println(a.color.PrintRedString(token))
	return token, nil
}

func (a *AirGap) InstallAirGapFiles(mirror string) error {
	fmt.Println(fmt.Sprintf("Downloading %s scripts and binaries", a.color.PrintBlueString("k3s")))
	install, err := a.embeddedFileLoader.LoadEmbeddedFile(installScript, common.TmpInstallScript)
	if err != nil {
		return err
	}
	err = a.remoteServer.TransferFile(install.Path, "/tmp/install.sh")
	if err != nil {
		return err
	}
	command, err := a.remoteServer.ExecuteCommand(common.Cmd1)
	if err != nil {
		return err
	}
	fmt.Println(command)

	if len(mirror) > 0 {
		registries = strings.Replace(registries, "repo", mirror, -1)
		reg, err := a.embeddedFileLoader.LoadEmbeddedFile(registries, common.TmpRegistriesYaml)
		if err != nil {
			return err
		}
		err = a.remoteServer.TransferFile(reg.Path, "/tmp/registries.yaml")
		if err != nil {
			return err
		}
		command, err := a.remoteServer.ExecuteCommand(common.InstallRegistriesYamlLocation)
		if err != nil {
			return err
		}
		fmt.Println(command)
	}

	binaryPath, err := a.airGapeFileDownloader.Download(a.base, a.binary)
	if err != nil {
		return err
	}
	err = a.remoteServer.TransferFile(binaryPath.Path, fmt.Sprintf("/tmp/%s", a.binary))
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf(common.Cmd2, fmt.Sprintf("/tmp/%s", a.binary), fmt.Sprintf("/tmp/%s", a.binary))
	executeCommand, err := a.remoteServer.ExecuteCommand(cmd)
	if err != nil {
		return err
	}
	fmt.Println(executeCommand)

	imagePath, err := a.airGapeFileDownloader.Download(a.base, a.images)
	err = a.remoteServer.TransferFile(imagePath.Path, fmt.Sprintf("/tmp/%s", a.images))
	if err != nil {
		return err
	}

	cmd = fmt.Sprintf(common.Cmd3, fmt.Sprintf("/tmp/%s", a.images))
	executeCommand, err = a.remoteServer.ExecuteCommand(cmd)
	if err != nil {
		return err
	}
	fmt.Println(executeCommand)

	return nil
}

func (a *AirGap) InstallControlPlaneNode() error {
	fmt.Println(fmt.Sprintf("Bootstraping %s cluster", a.color.PrintBlueString("k3s")))
	run, err := a.remoteServer.ExecuteCommand(fmt.Sprintf(common.InstallCmd, a.installK3sExec))
	if err != nil {
		return err
	}
	fmt.Println(run)
	return nil
}

func (a *AirGap) InstallWorkerNode(controlPlaneIp, token string) error {
	fmt.Println(fmt.Sprintf("Joining existing %s cluster %s\n", color.BlueString("k3s"), a.color.PrintGreenString(controlPlaneIp)))

	joinCMD := fmt.Sprintf(common.JoinCmd, token, controlPlaneIp)
	join, err := a.remoteServer.ExecuteCommand(joinCMD)
	if err != nil {
		return err
	}
	fmt.Println(join)

	return nil

	panic("implement me")
}

func (a *AirGap) GetKubeConfig() error {
	fmt.Println(fmt.Sprintf("Downloading %s kubeconfig ", a.color.PrintBlueString("k3s")))
	run, err := a.remoteServer.ExecuteCommand(common.KubeConfigCmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(run)
	run = strings.NewReplacer("localhost", a.remoteServer.GetRemoteServerIP(), "127.0.0.1", a.remoteServer.GetRemoteServerIP()).Replace(run)
	err = a.helper.WriteFile(common.K3sYaml, run)
	if err != nil {
		return err
	}
	return nil
}

func NewAirGap(base, arch, key, ip, user string, sudo bool) *AirGap {
	airGap := &AirGap{
		base:                  "https://github.com/k3s-io/k3s/releases/download/v1.21.1%2Bk3s1/",
		binary:                common.K3sBinary,
		airGapeFileDownloader: &downloader.AirGapeFileDownloader{},
		embeddedFileLoader:    &embedded.EmbeddedFileLoader{},
		remoteServer:          server.NewRemoteServer(key, ip, user, sudo),
		color:                 &term.Color{},
		helper:                &common.Helper{},
	}
	if arch == "amd64" {
		airGap.images = common.Amd64BinaryName
	} else {
		airGap.images = common.ArmBinaryName
	}
	if len(base) > 0 {
		airGap.base = base
	}
	return airGap
}
