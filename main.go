package main

import (
	"strconv"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hexbotio/hex-plugin"
	"golang.org/x/crypto/ssh"
)

type HexSsh struct {
}

func (g *HexSsh) Perform(args hexplugin.Arguments) (resp hexplugin.Response) {

	// initialize return values
	var output = ""
	var success = true

	// setup the server connection
	serverconn := true
	clientconn := &ssh.ClientConfig{
		User: args.Config["login"],
		Auth: []ssh.AuthMethod{
			ssh.Password(args.Config["pass"]),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	port := "22"
	if args.Config["port"] != "" {
		port = args.Config["port"]
	}
	limit, err := strconv.Atoi(args.Config["retry"])
	if err != nil {
		limit = 1
	}
	retryCounter := 0
	for retryCounter < limit {
		client, err := ssh.Dial("tcp", args.Config["server"]+":"+port, clientconn)
		if err != nil {
			output = err.Error()
			success = false
		}
		if client == nil {
			serverconn = false
		} else {
			defer client.Close()
			session, err := client.NewSession()
			if err != nil {
				output = err.Error()
				success = false
			}
			if session == nil {
				serverconn = false
			} else {
				defer session.Close()
				b, err := session.CombinedOutput(args.Command)
				output = string(b[:])
				if err != nil {
					output = output + "\n" + err.Error()
					success = false
				}
			}
		}
		if serverconn {
			retryCounter = limit
		} else {
			retryCounter += 1
			time.Sleep(time.Duration(3*retryCounter) * time.Second)
		}
	}
	if !serverconn {
		output = "ERROR - Cannot connect to server " + args.Config["server"]
		success = false
	}

	resp = hexplugin.Response{
		Output:  output,
		Success: success,
	}
	return resp
}

func main() {
	var pluginMap = map[string]plugin.Plugin{
		"action": &hexplugin.HexPlugin{Impl: &HexSsh{}},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: hexplugin.GetHandshakeConfig(),
		Plugins:         pluginMap,
	})
}
