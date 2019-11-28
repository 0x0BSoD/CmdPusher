package CmdPusher

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Cmd is a struct for command which we want to run
type Cmd struct {
	Commands   []string
	CurrentDir string
	StdOut     io.Writer
	StdErr     io.Writer
	ReturnCode int
}

type Client struct {
	Host     string
	Port     string
	User     string
	Password string
	AuthKey  string
	Insecure bool
	Timeout  time.Duration
	Session  *ssh.Session
}

func (c *Client) Connect() error {
	var checkKey ssh.HostKeyCallback
	var auth ssh.AuthMethod

	if c.Insecure {
		checkKey = ssh.InsecureIgnoreHostKey()
	} else {
		checkKey = ssh.FixedHostKey(getHostKey(c.Host))
	}

	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}

	if c.AuthKey != "" {
		key, err := ioutil.ReadFile(c.AuthKey)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}

		auth = ssh.PublicKeys(signer)
	} else {
		auth = ssh.Password(c.Password)
	}

	// ssh client config
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: checkKey,
		Timeout:         c.Timeout,
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), config)
	if err != nil {
		return err
	}

	// Create a session
	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	c.Session = sess

	return nil
}

func (c *Client) Close() error {
	return c.Session.Close()
}

func (c *Client) Run(cmd *Cmd) error {
	if c.Session == nil {
		return fmt.Errorf("session not started or already cloded")
	}

	c.Session.Stdout = cmd.StdOut
	c.Session.Stderr = cmd.StdErr

	stdin, err := c.Session.StdinPipe()
	if err != nil {
		return err
	}

	err = c.Session.Shell()
	if err != nil {
		return err
	}

	if cmd.CurrentDir != "" {
		_, err = fmt.Fprintf(stdin, "cd %s\n", cmd.CurrentDir)
		if err != nil {
			return err
		}
	}

	for _, cmd := range cmd.Commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(stdin, "%s\n", "exit $(echo $?)")
	if err != nil {
		return err
	}

	err = c.Session.Wait()
	if err != nil {
		log.Println("wait error")
		return err
	}

	return nil
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
