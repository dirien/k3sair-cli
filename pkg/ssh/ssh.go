package ssh

import (
	"bytes"
	"fmt"
	"github.com/k3sair/pkg/common"
	"github.com/morikuni/aec"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
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
	fmt.Println(fmt.Sprintf("Running remote command %s", aec.GreenB.Apply(cmd)))
	joinCMD := fmt.Sprintf(common.JoinCmdPart2, cmd, s.controlPlaneIP)
	run, err := s.RemoteRun(joinCMD, false)
	if err != nil {
		return "", err
	}
	return run, nil
}

func (s *SSH) RemoteRun(cmd string, join bool) (string, error) {
	fmt.Println(fmt.Sprintf("Running remote command %s", aec.LightMagentaF.Apply(cmd)))
	key, err := ioutil.ReadFile(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var ip string
	if join {
		ip = s.controlPlaneIP
	} else {
		ip = s.remoteIP
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(cmd)
	return strings.TrimSuffix(b.String(), "\n"), err
}

func (s *SSH) TransferFile(src *string, dstPath string) error {
	key, err := ioutil.ReadFile(s.privateKey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", s.remoteIP), config)
	if err != nil {
		return err
	}

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	dstFile, err := sftp.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := dstFile.Write([]byte(*src)); err != nil {
		return err
	}
	return nil
}
