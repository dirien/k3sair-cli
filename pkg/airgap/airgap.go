package airgap

import (
	_ "embed"
	"fmt"
	"github.com/k3sair/pkg/common"
	"github.com/k3sair/pkg/ssh"
	"github.com/morikuni/aec"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AirGap struct {
	base   string
	binary string
	images string
	key    string
	ssh    *ssh.SSH
	sudo   bool
	user   string
}

//go:embed install.sh
var installScript string

type AirGapped interface {
	DownloadAirGap() error
	Install() error
	Join() error
}

func (a *AirGap) Join() error {
	command := common.CheckSudo(a.sudo, common.Cmd4)
	token, err := a.ssh.RemoteRun(command, true)
	if err != nil {
		return err
	}
	token = strings.TrimSuffix(token, "\n")
	fmt.Println(token)

	joinCMD := fmt.Sprintf(common.JoinCmd, token)
	join, err := a.ssh.RemoteJoinRun(joinCMD)
	if err != nil {
		return err
	}
	fmt.Println(join)

	return nil
}

func (a *AirGap) Install() error {
	run, err := a.ssh.RemoteRun(common.InstallCmd, false)
	if err != nil {
		return err
	}
	fmt.Println(run)
	return nil
}

func (a *AirGap) DownloadAirGap() error {
	binaryPath, err := download(a.base, a.binary)
	if err != nil {
		return err
	}
	err = copyBinary(binaryPath, a)
	if err != nil {
		return err
	}
	imagePath, err := download(a.base, a.images)
	if err != nil {
		return err
	}
	err = copyImage(imagePath, a)
	if err != nil {
		return err
	}
	return nil
}

func copyBinary(path string, a *AirGap) error {
	tmpFolder, err := transfer(path, a.ssh, a.binary)
	if err != nil {
		return err
	}
	command := common.CheckSudo(a.sudo, common.Cmd2)
	cmd := fmt.Sprintf(command, tmpFolder, tmpFolder)
	err = runRemoteCmd(err, a, cmd)
	if err != nil {
		return err
	}
	return nil
}

func runRemoteCmd(err error, a *AirGap, cmd string) error {
	run, err := a.ssh.RemoteRun(cmd, false)
	if err != nil {
		return err
	}
	fmt.Println(run)
	return nil
}

func transfer(path string, ssh *ssh.SSH, binary string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	bstring := string(b)
	tmpFolder := fmt.Sprintf("/tmp/%s", binary)
	fmt.Println(fmt.Sprintf("Start transfer file %s to remote server", aec.RedF.Apply(tmpFolder)))
	err = ssh.TransferFile(
		&bstring,
		tmpFolder)
	if err != nil {
		return "", err
	}
	return tmpFolder, nil
}

func copyImage(path string, a *AirGap) error {
	tmpFolder, err := transfer(path, a.ssh, a.images)
	if err != nil {
		return err
	}
	command := common.CheckSudo(a.sudo, common.Cmd3)
	cmd := fmt.Sprintf(command, tmpFolder)
	err = runRemoteCmd(err, a, cmd)
	if err != nil {
		return err
	}
	return nil
}

func download(base, file string) (path string, err error) {
	fmt.Println(fmt.Sprintf("Download Air-Gap file %s", aec.GreenF.Apply(file)))
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	p := filepath.FromSlash(fmt.Sprintf("%s/%s", tmp, file))
	out, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer out.Close()

	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	c := &http.Client{Transport: transport}

	resp, err := c.Get(fmt.Sprintf("%s/%s", base, file))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(fmt.Sprintf("Air-Gap file succesfully downloaded at %s", aec.LightRedF.Apply(p)))
	return p, nil
}

func NewAirGap(base, arch, key, ip, controlPlaneIp, user string, sudo bool) *AirGap {
	ssh := ssh.NewAirGapOperations(key, ip, controlPlaneIp, user)
	err := ssh.TransferFile(&installScript, fmt.Sprintf(common.InstallScriptLocation, user))
	if err != nil {
		return nil
	}

	command := common.CheckSudo(sudo, common.Cmd1)
	run, err := ssh.RemoteRun(command, false)
	if err != nil {
		return nil
	}
	fmt.Println(run)
	ptr := &AirGap{
		base:   "https://github.com/k3s-io/k3s/releases/download/v1.21.1%2Bk3s1/",
		binary: common.K3sBinary,
		key:    key,
		ssh:    ssh,
		sudo:   sudo,
		user:   user,
	}
	if arch == "amd64" {
		ptr.images = common.Amd64BinaryName
	} else {
		ptr.images = common.ArmBinaryName
	}
	if len(base) > 0 {
		ptr.base = base
	}
	return ptr
}
