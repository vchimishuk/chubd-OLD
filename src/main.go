package main

import (
	"os"
	"fmt"
	"os/signal"
	"./server"
	"./protocol"
)

// UNIX signals
const (
	SigTerm = 15
)

func main() {
	// TODO: Parse command line parameters.

	// TODO: Daemonize itself.
	//       daemonize()

	host := "127.0.0.1"
	port := 8888

	srv, err := server.NewTCPServer(host, port)
	if err != nil {
		fmt.Printf("Failed to create server. %s", err)
		os.Exit(1)
	}

	// Run listening loop.
	srv.SetConnectionHandler(new(protocol.ConnectionHandler))
	go srv.Serve()

	// On SIGTERM received we have to close all client connections
	// and then exit. So we loop till this signal will be recieved.
	for {
		sig := (<-signal.Incoming).(os.UnixSignal) // XXX: Will works in windows? No way!
		sigNum := int32(sig)

		if sigNum == SigTerm {
			break
		}
	}
}
