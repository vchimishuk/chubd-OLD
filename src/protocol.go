// Protocol package implements client communication protocol.
package protocol

import (
	"os"
	"net"
	"bufio"
	"fmt"
	"scanner"
	"strings"
	"./server"
	"./filesystem"
)

// command represents parsed command.
type command struct {
	// Name of the command.
	Name string
	// Command parameters represented as strings (as they were sent by client).
	Parameters []string
}

// parseCommand parses client's string command (request) to command object.
// TODO: This parser sucks, so write you own with blackjac & bitches.
func parseCommand(str string) (cmd *command, err os.Error) {
	cmd = new(command)

	var s scanner.Scanner
	s.Init(strings.NewReader(str))

	// Command name.
	tok := s.Scan()
	if tok != scanner.EOF {
		cmd.Name = strings.Trim(s.TokenText(), "\"")
	}

	// Parameters.
	cmd.Parameters = make([]string, 0)
	tok = s.Scan()
	for tok != scanner.EOF {
		cmd.Parameters = append(cmd.Parameters, strings.Trim(s.TokenText(), "\""))
		tok = s.Scan()
	}

	return cmd, nil
}

// CommandHandler signature.
type commandHandler func(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error

// commandDescriptor describes command attribetes (parameters, handler function, ...)
type commandDescriptor struct {
	// Number of required arguments.
	argc int
	// Handler function for the command.
	handler commandHandler
}

// All supported commands descriptors. 
var commandDescriptors = map[string]commandDescriptor{
	"CD":   commandDescriptor{1, cmdCd},
	"LS":   commandDescriptor{0, cmdLs},
	"PING": commandDescriptor{0, cmdPing},
	"PWD":  commandDescriptor{0, cmdPwd},
	// "QUIT": built-in
}

// CommandHandler struct.
type CommandHandler struct {
	fs *filesystem.Filesystem
}

// NewCommandHandler creates new initialized command handler object.
func NewCommandHandler() *CommandHandler {
	ch := new(CommandHandler)
	ch.fs = filesystem.New()

	return ch
}

// HandleCommand interface implementation which will be called on every client's request
func (ch *CommandHandler) HandleCommand(writer *bufio.Writer, request string) bool {
	writeOk := func() {
		writer.WriteString("OK\n")
		writer.Flush()
	}

	writeError := func(err os.Error) {
		writer.WriteString(fmt.Sprintf("ERROR %s\n", err))
		writer.Flush()
	}

	cmd, err := parseCommand(request)
	if err != nil {
		writeError(err)
		return false
	}

	// Check if it is built-in QUIT command.
	if cmd.Name == "QUIT" {
		writeOk()
		return true
	}

	err = ch.run(writer, cmd)
	if err == nil {
		writeOk()
	} else {
		writeError(err)
	}

	return false
}

// run multiplex all handlers and select related to be invoked.
func (ch *CommandHandler) run(writer *bufio.Writer, cmd *command) os.Error {
	// Check if command is supported.
	cmdDescriptor, ok := commandDescriptors[cmd.Name]
	if !ok {
		return os.NewError(fmt.Sprintf("Unsupported command '%s'", cmd.Name))
	}

	// Check if number of parameters are correct.
	if cmdDescriptor.argc != len(cmd.Parameters) {
		return os.NewError(fmt.Sprintf("Command '%s' requires %d parameters but %d given",
			cmd.Name, cmdDescriptor.argc, len(cmd.Parameters)))
	}

	return cmdDescriptor.handler(ch, writer, cmd)
}

// ConnectionHandler implementation type.
type ConnectionHandler func(conn net.Conn) CommandHandler

// HandleConnection handle every client's connection.
func (ch ConnectionHandler) HandleConnection(conn net.Conn) server.CommandHandler {
	return NewCommandHandler()
}

// cmdPing implements PING command.
// PING command just does nothing.
func cmdPing(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	return nil
}

// cmdPwd implemets PWD command.
// PWD command prints current working directory.
func cmdPwd(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	writer.WriteString(ch.fs.WorkingDir())
	writer.WriteString("\n")

	return nil
}

// cmdCd implements LS server command.
// CD changes current working directory.
// Parameters:
// * directory
func cmdCd(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	dir := cmd.Parameters[0]

	return ch.fs.SetWorkingDir(dir)

}

// cmdLs implements LS server command.
// LS command prints sorted (dirs before files) working direcory listing.
func cmdLs(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	// Write HTTP header-like string to writer.
	// Key: Value
	writef := func(writer *bufio.Writer, format string, args ...interface{}) {
		str := fmt.Sprintf(format, args...)
		writer.WriteString(str)
	}

	entries, err := ch.fs.List()
	if err != nil {
		return err
	}

	lastIndex := len(entries) - 1
	//for i, entry := range entries {
	for i := 0; i < len(entries); i++ {
		writef(writer, "Type: %s\n", entries[i].TypeString())

		switch entries[i].Type() {
		case filesystem.TypeTrack:
			writef(writer, "FileName: %s\n", entries[i].Track().Filename)
		case filesystem.TypeDirectory:
			writef(writer, "Name: [%s]\n", entries[i].Directory().Name)
		}

		if i < lastIndex {
			writer.WriteString("\n")
		}
	}

	return nil
}
