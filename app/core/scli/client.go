package scli

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"pear-admin-go/app/core/log"

	"pear-admin-go/app/model"
	"time"
)

var client *sftp.Client

func Instance(server model.TaskServer) (*sftp.Client, error) {
	//if client != nil {
	//	return client, nil
	//}
	var err error
	if server.ConnType == 1 {
		client, err = connectByPwd(server)
	} else {
		client, err = connectByKey(server)
	}
	if err != nil {
		log.Instance().Error(err.Error())
		return nil, err
	}
	return client, nil
}

func connectByPwd(server model.TaskServer) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(server.ServerPassword))

	clientConfig = &ssh.ClientConfig{
		User:    server.ServerAccount,
		Auth:    auth,
		Timeout: 15 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", server.ServerIp, server.Port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		log.Instance().Error("ssh.Dial error: " + err.Error())
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Instance().Error("sftp.NewClient error: " + err.Error())
		return nil, err
	}
	return sftpClient, nil
}

func connectByKey(server model.TaskServer) (*sftp.Client, error) {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	)
	key, err := ioutil.ReadFile(server.PrivateKeySrc)
	if err != nil {
		log.Instance().Error("unable to read private key: " + err.Error())
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Instance().Error("unable to parse private key: " + err.Error())
		return nil, err
	}

	clientConf := &ssh.ClientConfig{
		User: server.ServerAccount,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 15 * time.Second,
	}
	sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.ServerIp, server.Port), clientConf)
	if err != nil {
		log.Instance().Error("unable to connect: " + err.Error())
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Instance().Error("sftp.NewClient error: " + err.Error())
		return nil, err
	}
	return sftpClient, nil
}
