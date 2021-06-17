package ssh

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/common"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
)

type SSH struct {
	remoteIP       string
	privateKey     string
	controlPlaneIP string
	user           string
}

type AirGapOperations interface {
	TransferFile(src, dstPath string) error
	RemoteRun(cmd string, join bool) (string, error)
	RemoteJoinRun(cmd string) (string, error)
}

func NewAirGapOperations(privateKey, ip, controlPlaneIp, user string) *SSH {
	ssh := &SSH{
		remoteIP:       ip,
		controlPlaneIP: controlPlaneIp,
		privateKey:     privateKey,
		user:           user,
	}
	return ssh
}

func (s *SSH) RemoteJoinRun(cmd string) (string, error) {
	joinCMD := fmt.Sprintf(common.JoinCmdPart2, cmd, s.controlPlaneIP)
	run, err := s.RemoteRun(joinCMD, false)
	if err != nil {
		return "", err
	}
	return run, nil
}

func (s *SSH) RemoteRun(cmd string, join bool) (string, error) {
	fmt.Println(fmt.Sprintf("Running remote command %s", color.GreenString(cmd)))

	var ip string
	if join {
		ip = s.controlPlaneIP
	} else {
		ip = s.remoteIP
	}

	auth, err := goph.Key(s.privateKey, "")
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User:     s.user,
		Addr:     ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	out, err := client.Run(cmd)
	return string(out), err
}

func (s *SSH) TransferFile(path *string, dstPath string) error {
	auth, err := goph.Key(s.privateKey, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     s.user,
		Addr:     s.remoteIP,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	err = client.Upload(*path, dstPath)
	if err != nil {
		return err
	}

	return nil
}
