// Protocol package implements client communication protocol.
package protocol

import (
	"os"
	"net"
	"bufio"
	"fmt"
	"scanner"
	"strings"
	"strconv"
	"./server"
	"./vfs"
	"./player"
)

// command represents parsed command.
type command struct {
	// Name of the command.
	Name string
	// Command parameters represented as strings (as they were sent by client).
	Parameters []string
}

// Field constants.
const (
	fieldNameType     = "Type"
	fieldNameFilename = "Filename"
	fieldNameArtist   = "Artist"
	fieldNameAlbum    = "Album"
	fieldNameTitle    = "Title"
	fieldNameNumber   = "Number"
	fieldNameLength   = "Length"
	fieldNameName     = "Name"
)

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
	"CD":             commandDescriptor{1, cmdCd},
	"LS":             commandDescriptor{0, cmdLs},
	"PING":           commandDescriptor{0, cmdPing},
	"PWD":            commandDescriptor{0, cmdPwd},
	"PLAYLISTS":      commandDescriptor{0, cmdPlaylists},
	"ADDPLAYLIST":    commandDescriptor{1, cmdAddPlaylist},
	"DELETEPLAYLIST": commandDescriptor{1, cmdDeletePlaylist},
	"PLAYVFS":        commandDescriptor{1, cmdPlayVfs},
	"PAUSE":          commandDescriptor{0, cmdPause},
	"KILL":           commandDescriptor{0, cmdKill},
	// "QUIT": built-in
}

// CommandHandler struct.
type CommandHandler struct {
	fs *vfs.Filesystem
}

// NewCommandHandler creates new initialized command handler object.
func NewCommandHandler() *CommandHandler {
	ch := new(CommandHandler)
	ch.fs = vfs.New()

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
	writePair := func(key string, value string) {
		writer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	entries, err := ch.fs.List()
	if err != nil {
		return err
	}

	lastIndex := len(entries) - 1
	for i := 0; i < len(entries); i++ {
		writePair(fieldNameType, entries[i].TypeString())

		switch entries[i].Type() {
		case vfs.TypeTrack:
			track := entries[i].Track()
			tag := track.Tag
			// Tracks are indentified by filename:trackNum scheme.
			// For single track files trackNum is 0.
			writePair(fieldNameFilename, fmt.Sprintf("%s:%d", track.FilePath.Path(), track.Number))
			writePair(fieldNameArtist, tag.Artist)
			writePair(fieldNameAlbum, tag.Album)
			writePair(fieldNameTitle, tag.Title)
			writePair(fieldNameNumber, tag.Number)
			writePair(fieldNameLength, tag.Length)
		case vfs.TypeDirectory:
			dir := entries[i].Directory()
			writePair(fieldNameFilename, dir.Filename.Path())
			writePair(fieldNameName, dir.Name)
		}

		if i < lastIndex {
			writer.WriteString("\n")
		}
	}

	return nil
}

// cmdPlaylists handles PLAYLISTS server command.
// PLAYLISTS command prints list of the registered playlists.
func cmdPlaylists(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	lastIndex := len(player.Playlists()) - 1

	for i, pl := range player.Playlists() {
		// System playlists are not visible to clients.
		if pl.IsSystem() {
			continue
		}

		writer.WriteString(fmt.Sprintf("Name: %s\n", pl.Name()))
		writer.WriteString(fmt.Sprintf("Length: %d\n", pl.Len()))

		if i < lastIndex {
			writer.WriteString("\n")
		}
	}

	return nil
}

// cmdAddPlaylist creates new empty playlist.
// Parameters:
// * playlist name
func cmdAddPlaylist(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	name := cmd.Parameters[0]

	return player.AddPlaylist(name)
}

// cmdDeletePlaylist deletes existing playlist for the given name.
// Parameters:
// * playlist name
func cmdDeletePlaylist(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	name := cmd.Parameters[0]

	return player.DeletePlaylist(name)
}

// cmdPlayVfs plays track from the working directory.
// Parameters:
// * filename in next format: file.flac:3
func cmdPlayVfs(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	param := cmd.Parameters[0]
	i := strings.LastIndex(param, ":")
	if i == -1 {
		return os.NewError("Bad filename format. Expected filename:track_number")
	}

	filename := param[:i]
	trackNumber, err := strconv.Atoi(param[i + 1:])
	fmt.Printf("trackNumber: %d\n\n", trackNumber)
	if err != nil || trackNumber < 0 {
		return os.NewError("Bad track number format. 0..99 expected")
	}

	pl, _ := player.Playlist(vfs.PlaylistName) // I don't check error, because this playlist should be present always.
	pl.Clear()

	entries, err := ch.fs.List()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Type() == vfs.TypeTrack {
			pl.Append(entry.Track())
		}
	}

	// Find track position in the VFS playlist.
	pos := -1
	for i, track := range pl.Tracks() {
		if track.Number == trackNumber && track.FilePath.Path() == filename {
			pos = i
		}
	}
	if pos == -1 {
		return os.NewError("Track not found in working directory")
	}
	
	player.Play(vfs.PlaylistName, pos)

	return nil
}

// cmdPause toggle player's pause state.
func cmdPause(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	player.Pause()

	return nil
}

// cmdKill stops player. After program can be terminated.
func cmdKill(ch *CommandHandler, writer *bufio.Writer, cmd *command) os.Error {
	player.Stop()

	// TOOD: Add here correct program termination.
	panic("TOOD: Add here correct program termination")

	return nil
}
