package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func DialWithKey(addr, user, keyfile string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fmt.Println("Inside the hostkey call back")
			return nil
		}),
	}

	return Dial("tcp", addr, config)
}
func Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func ExecuteCommand(client *ssh.Client, cmd string) {
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("\n Error in creating session %+v", err)
		return
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		fmt.Printf("Error in executing command %+v\n", err)
		return
	}
	fmt.Printf("%s", stdoutBuf.String())

}
func ExecuteMultipleCommand(client *ssh.Client, commands []string) {
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("\n Error in creating session %+v", err)
		return
	}
	defer session.Close()
	cmdInputPipe, err := session.StdinPipe()
	if err != nil {
		fmt.Printf("\n Error in opening input pipe %+v", err)
		return
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err != nil {
		fmt.Printf("\n Error in opening output pipe %+v", err)
		return
	}
	err = session.Shell()
	if err != nil {
		fmt.Printf("\n Error in opening shell %+v", err)
		return
	}
	for _, cmd := range commands {
		cmdStr := fmt.Sprintf("%s\n", cmd)
		cmdInputPipe.Write([]byte(cmdStr))
	}
	cmdInputPipe.Write([]byte("exit\n"))
	err = session.Wait()
	if err != nil {
		fmt.Printf("\n Error in opening shell %+v", err)
	}
}
func main() {
	client, err := DialWithKey("35.231.83.48:22", "purnendu", "/home/suddutt1/keys/pjm-google.pem")
	if err != nil {
		fmt.Printf("Error in opening connection %+v\n", err)
		return
	}
	defer client.Close()
	//ExecuteCommand(client, "ls -ltr")
	ExecuteMultipleCommand(client, []string{"cd simplemv/common", "source setFabricEnv.sh ", "cd ../orderer", "docker-compose up -d", "docker ps"})
}
