package server

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/k3sair/pkg/common"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

type Operations interface {
	TransferFile(src, dstPath string) error
	ExecuteCommand(cmd string) (string, error)
	GetRemoteServerIP() string
}

type RemoteServer struct {
	ip            string
	port          uint
	privateSSHKey string
	user          string
	sudo          bool
}

func NewRemoteServer(privateKey, ip, user string, port uint, sudo bool) *RemoteServer {
	ssh := &RemoteServer{
		ip:            ip,
		port:          port,
		privateSSHKey: privateKey,
		user:          user,
		sudo:          sudo,
	}
	return ssh
}

func (r *RemoteServer) GetRemoteServerIP() string {
	return r.ip
}

func (r *RemoteServer) TransferFile(src, dstPath string) error {
	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		return err
	}

	client, err := goph.NewConn(&goph.Config{
		User:     r.user,
		Addr:     r.ip,
		Port:     r.port,
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
	fmt.Printf("Running remote command %s\n", color.GreenString(cmd))

	auth, err := goph.Key(r.privateSSHKey, "")
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User:     r.user,
		Addr:     r.ip,
		Port:     r.port,
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	command := common.CheckSudo(r.sudo, cmd)
	out, err := client.Run(command)
	return string(out), err
}
