# CmdPusher

Simple wrapper around crypto/ssh package


```go
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

cmd := &Cmd{
    Commands:   []string{"ls -la", "echo foo"},
    CurrentDir: "/home",
    StdOut:     os.Stdout,
    StdErr:     os.Stderr,
}

err := CLIENT.Connect()
if err != nil {
    panic(err)
}

err = CLIENT.Run(cmd)
if err != nil {
    panic(err)
}

_ = CLIENT.Close()

```