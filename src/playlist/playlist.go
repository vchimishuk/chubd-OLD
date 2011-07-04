// Playlist package implements playlist, -- ordered list of tracks to be played, abstraction.
package playlist

import (
	"sync"
	"./vfs"
)

// Playlist
type Playlist struct {
	// Uniq name of the playlist.
	// Playlists with name matched with '\*.*\*' mask are system and can't be
	// created or deleted by user.
	name string
	// List of tracks present in the current playlist.
	tracks []*vfs.Track
	// Lock this mutex for operations which changes playlist.
	mtx sync.Mutex
}

// New returns an initialized playlist.
func New(name string) *Playlist {
	pl := new(Playlist)
	pl.name = name

	return pl
}

// Name returns name of the playlist.
func (pl *Playlist) Name() string {
	return pl.name
}

// Tracks returns list of tracks presented in the playlist.
func (pl *Playlist) Tracks() []*vfs.Track {
	return pl.tracks
}

// Track returns track by its position.
func (pl *Playlist) Track(n int) *vfs.Track {
	return pl.tracks[n]
}

// Len returns the total number of tracks present in playlist. 
func (pl *Playlist) Len() int {
	return len(pl.tracks)
}

// IsSystem returns true if playlist is system and can't be created/deleted by user.
// System playlist is a playlist that has a name strats and ends with asterisk.
func (pl *Playlist) IsSystem() bool {
	return len(pl.name) >= 2 && pl.name[0] == '*' && pl.name[len(pl.name)-1] == '*'
}

// Append appends new tracks at the end of the playlist.
func (pl *Playlist) Append(tracks ...*vfs.Track) {
	pl.mtx.Lock()
	defer pl.mtx.Unlock()

	pl.tracks = append(pl.tracks, tracks...)
}

// Clear removes all items from the playlist.
func (pl *Playlist) Clear() {
	pl.mtx.Lock()
	defer pl.mtx.Unlock()

	pl.tracks = make([]*vfs.Track, 0)
}
