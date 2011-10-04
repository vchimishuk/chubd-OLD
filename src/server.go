// Server implements clients handling functionality.
package server

import (
	"os"
	"net"
	"net/textproto"
	"bufio"
	"fmt"
	"sync"
)

// Server is the server itself interface, which can listen network (TCP, UNIX Sockets, etc.)
// connections and process client's commands with the CommandHandler interface implementations.
type Server interface {
	// SetConnectionHandler should be called before Serve method
	// to set client connection processor.
	SetConnectionHandler(handler ConnectionHandler)
	// Serve start main server loop, which accepts connection form the clients
	// and do next communication processing.
	Serve() os.Error
}

// HandleConnection method is called every time client is connected. And this method should returns
// implementation of the CommandHandler interface, which will handle every client's command.
type ConnectionHandler interface {
	// HandleConnection calls every time new client is connected.
	// If HandleConnection return nil client should be rejected and connection closed.
	HandleConnection(conn net.Conn) CommandHandler
}

// CommandHandler is the client requests handler. Every time client send some command
// to server HandleCommand method will be called.
type CommandHandler interface {
	// HandleCommand calls every time new command from client recived.
	// true return result means that communication was ended. Than server close connection.
	HandleCommand(writer *bufio.Writer, command string) bool
}

// tcpServer represents server which works on TCP/IP netwoks.
type tcpServer struct {
	// Number of currently connected clients.
	clientsCount int
	// Mutex for protecting clientsCount field.
	mutex             sync.Mutex
	listener          *net.TCPListener
	connectionHandler ConnectionHandler
	commandHandler    CommandHandler
}

// NewTCPServer creates newly initialized TCP server.
func NewTCPServer(ip string, port int) (srv Server, err os.Error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		msg := fmt.Sprintf("'%s' is not correct IP address.", ip)
		return nil, os.NewError(msg)
	}
	addr := net.TCPAddr{ipAddr, port}
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		return nil, err
	}

	server := new(tcpServer)
	server.listener = listener

	return server, nil
}

// Set connection handler.
func (srv *tcpServer) SetConnectionHandler(handler ConnectionHandler) {
	srv.connectionHandler = handler
}

// Serve starts TCP server main loop and make it ready to clients handling.
func (srv *tcpServer) Serve() os.Error {
	if srv.connectionHandler == nil {
		return os.NewError("SetConnectionHandler should be call first")
	}

	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			return err
		}

		commandHandler := srv.connectionHandler.HandleConnection(conn)
		if commandHandler != nil {
			go srv.handleClient(conn, commandHandler)
		} else {
			conn.Close()
		}
	}

	return nil
}

// helloMessage returns server greeting string.
func (srv *tcpServer) helloMessage() string {
	return "Chubd 0.0 service\nOK\n" // TODO: Retrive version number from some global spot.
}

// handleClient handle the client and organize communication between the client
// and CommandHandler.
func (srv *tcpServer) handleClient(conn net.Conn, commandHandler CommandHandler) {
	srv.addClient()
	defer srv.removeClient()

	// Say hello to client.
	writer := bufio.NewWriter(conn)
	writer.WriteString(srv.helloMessage())
	writer.Flush()

	for {
		reader := textproto.NewReader(bufio.NewReader(conn))
		command, err := reader.ReadLine() // TODO: Parse request string to command.
		if err != nil {
			return // Connection was closed by client, or something like that.
		}

		exit := commandHandler.HandleCommand(bufio.NewWriter(conn), command)
		if exit {
			break // Client wants to end this conversation.
		}
	}

	conn.Close()
}

// addClient register new client.
func (srv *tcpServer) addClient() {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.clientsCount++
}

// removeClient delete new client.
func (srv *tcpServer) removeClient() {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.clientsCount--
}
