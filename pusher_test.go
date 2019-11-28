package CmdPusher

import (
	"log"
	"testing"
)

func TestClient_Connect_Password(t *testing.T) {
	c := Client{
		Host:     "you_server",
		Port:     "22",
		User:     "testuser",
		Password: "test",
		AuthKey:  "",
		Insecure: true,
		Timeout:  0,
		Session:  nil,
	}

	err := c.Connect()
	if err != nil {
		panic(err)
	}

	log.Println("Connected!")

	err = c.Close()
	if err != nil {
		panic(err)
	}
	log.Println("Disconnected!")
}
