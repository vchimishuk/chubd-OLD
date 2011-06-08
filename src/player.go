// player package is the core of the program: it manages playlists and player's state.
package player

import (
	"os"
	"fmt"
	"sync"
	"./playlist"
)

var (
	// All (user and system) playlists list.
	playlists []*playlist.Playlist
	// Mutex for protecting playlists modification operations.
	playlistsLock sync.Mutex
)


// Playlists returns list of all existent playlists.
func Playlists() []*playlist.Playlist {
	return playlists
}

// Playlist returns playlist object by name.
func Playlist(name string) (playlist *playlist.Playlist, err os.Error) {
	if pl := getPlaylistByName(name); pl != nil {
		return pl, nil
	}

	return nil, os.NewError(fmt.Sprintf("Playlist '%s' not found"))
}

// AddPlaylist creates new empty playlist and adds it to the list.
func AddPlaylist(name string) os.Error {
	playlistsLock.Lock()
	defer playlistsLock.Unlock()

	// Check if name already exists.
	if pl := getPlaylistByName(name); pl != nil {
		return os.NewError(fmt.Sprintf("Playlist '%s' already exists", name))
	}

	pl := playlist.New(name)
	if pl.IsSystem() {
		return os.NewError(fmt.Sprintf("System playlist can't be created"))
	}

	playlists = append(playlists, pl)

	return nil
}

// DeletePlaylist deletes existing playlist by name.
func DeletePlaylist(name string) os.Error {
	playlistsLock.Lock()
	defer playlistsLock.Unlock()

	newPlaylists := make([]*playlist.Playlist, 0, len(playlists))
	for _, pl := range playlists {
		if pl.Name() != name {
			newPlaylists = append(newPlaylists, pl)
		} else if pl.IsSystem() {
			return os.NewError("System playlist can't be deleted")
		}
	}

	if len(newPlaylists) == len(playlists) {
		return os.NewError(fmt.Sprintf("Playlist '%s' not found", name))
	}

	playlists = newPlaylists

	return nil
}

// getPlaylistByName returns playlist for given name
// or nil if there is no such playlist registered.
func getPlaylistByName(name string) *playlist.Playlist {
	for _, pl := range playlists {
		if pl.Name() == name {
			return pl
		}
	}

	return nil
}

// Package init function.
func init() {
	playlists = make([]*playlist.Playlist, 0)

	// We have one system (predefined) playlist, -- *vfs*.
	playlists = append(playlists, playlist.New("*vfs*"))
}
