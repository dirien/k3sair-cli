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
	ExecuteCommand(cmd string) (string, error)
	GetRemoteServerIP() string
}

type RemoteServer struct {
	ip            string
	port          uint
	privateSSHKey string
	user          string
	sudo          bool
	helper        *common.Helper
}

func NewRemoteServer(privateKey, ip, user string, port uint, sudo bool) *RemoteServer {
	ssh := &RemoteServer{
		ip:            ip,
		port:          port,
		privateSSHKey: privateKey,
		user:          user,
		sudo:          sudo,
		helper:        &common.Helper{},
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
	fmt.Println(fmt.Sprintf("Running remote command %s", color.GreenString(cmd)))

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
	command := r.helper.CheckSudo(r.sudo, cmd)
	out, err := client.Run(command)
	return string(out), err
}
