package main

import (
	"os"
	"fmt"
	"./server"
	"./protocol"
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

	srv.SetConnectionHandler(new(protocol.ConnectionHandler))
	srv.Serve()
}
