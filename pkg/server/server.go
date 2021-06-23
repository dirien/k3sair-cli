package server

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/k3sair/pkg/common"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
)

type ServerOperations interface {
	TransferFile(src, dstPath string) error
	ExecuteCommand(cmd string, join bool) (string, error)
}

type RemoteServer struct {
	ip            string
	privateSSHKey string
	User          string
	sudo          bool
	helper        *common.Helper
}

func NewRemoteServer(privateKey, ip, user string, sudo bool) *RemoteServer {
	ssh := &RemoteServer{
		ip:            ip,
		privateSSHKey: privateKey,
		User:          user,
		sudo:          sudo,
		helper:        &common.Helper{},
	}
	return ssh
}

func (r *RemoteServer) TransferFile(src, dstPath string) error {
	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		return err
	}

	client, err := goph.NewConn(&goph.Config{
		User:     r.User,
		Addr:     r.ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}

	defer client.Close()
	err = client.Upload(src, dstPath)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemoteServer) ExecuteCommand(cmd string) (string, error) {
	fmt.Println(fmt.Sprintf("Running remote command %s", color.GreenString(cmd)))

	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User:     r.User,
		Addr:     r.ip,
		Port:     22,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	command := r.helper.CheckSudo(r.sudo, cmd)
	out, err := client.Run(command)
	return string(out), err
}

/*
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
*/
