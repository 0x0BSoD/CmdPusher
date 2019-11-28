package CmdPusher

import (
	"log"
	"os"
	"testing"
)

var CLIENT = Client{
	Host:     "you_server",
	Port:     "22",
	User:     "user",
	Password: "pass",
	AuthKey:  "",
	Insecure: true,
	Timeout:  0,
	Session:  nil,
}

func TestClient_Connect_Password(t *testing.T) {
	err := CLIENT.Connect()
	if err != nil {
		panic(err)
	}
	log.Println("Connected!")

	err = CLIENT.Close()
	if err != nil {
		panic(err)
	}
	log.Println("Disconnected!")
}

func TestClient_Run(t *testing.T) {
	err := CLIENT.Connect()
	if err != nil {
		panic(err)
	}
	log.Println("Connected!")

	cmd := &Cmd{
		Commands:   []string{"ls -la"},
		CurrentDir: "/etc/puppet/environments",
		StdOut:     os.Stdout,
		StdErr:     os.Stderr,
	}

	err = CLIENT.Run(cmd)
	if err != nil {
		panic(err)
	}

	_ = CLIENT.Close()

	log.Println("Disconnected!")
}