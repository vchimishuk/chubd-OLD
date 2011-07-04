// player package is the core of the program: it manages playlists and player's state.
package player

import (
	"os"
	"fmt"
	"sync"
	"./vfs"
	"./playlist"
	"./audio"
	"./ogg"
	"./alsa"
)

// Player mutex. All public player commands should be protected with this mutex lock.
var mutex sync.Mutex
// All (user and system) playlists list.
var playlists []*playlist.Playlist
// thread is the main player thread (goroutine wrapper).
var thread *playingThread

// Playlists returns list of all existent playlists.
func Playlists() []*playlist.Playlist {
	return playlists
}

// Playlist returns playlist object by name.
func Playlist(name string) (playlist *playlist.Playlist, err os.Error) {
	mutex.Lock()
	defer mutex.Unlock()

	playlist, err = getPlaylistByName(name)

	return
}

// AddPlaylist creates new empty playlist and adds it to the list.
func AddPlaylist(name string) os.Error {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if name already exists.
	pl, _ := getPlaylistByName(name)
	if pl != nil {
		return os.NewError(fmt.Sprintf("Playlist '%s' already exists", name))
	}

	pl = playlist.New(name)
	if pl.IsSystem() {
		return os.NewError(fmt.Sprintf("System playlist can't be created"))
	}

	playlists = append(playlists, pl)

	return nil
}

// DeletePlaylist deletes existing playlist by name.
func DeletePlaylist(name string) os.Error {
	mutex.Lock()
	defer mutex.Unlock()

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

// Play start playing existing a track form an existing playlist.
func Play(playlistName string, trackNumber int) os.Error {
	mutex.Lock()
	defer mutex.Unlock()

	pl, err := getPlaylistByName(playlistName)
	if err != nil {
		return err
	}

	if trackNumber < 0 || trackNumber >= pl.Len() {
		return os.NewError(fmt.Sprintf("Playlist '%s' has no track number %d.",
			playlistName, trackNumber))
	}

	thread.Play(pl.Track(trackNumber))

	return nil
}

// Pause pause or unpause playing process.
func Pause() {
	thread.Pause()
}

// getPlaylistByName returns playlist for given name
// or nil if there is no such playlist registered.
func getPlaylistByName(name string) (playlist *playlist.Playlist, err os.Error) {
	for _, pl := range playlists {
		if pl.Name() == name {
			return pl, nil
		}
	}

	return nil, os.NewError(fmt.Sprintf("Playlist '%s' not found"))
}

// Package init function.
func init() {
	// Audio tagreaders.
	audio.RegisterTagReaderFactory(ogg.NewTagReader)
	// Audio outputs.
	audio.RegisterOutput(alsa.DriverName, alsa.New)

	// Playists
	playlists = make([]*playlist.Playlist, 0)
	// We have one system (predefined) playlist, -- *vfs*.
	playlists = append(playlists, playlist.New(vfs.PlaylistName))

	// Create and start playing thread.
	thread = newPlayingThread()
	thread.Start()
}
