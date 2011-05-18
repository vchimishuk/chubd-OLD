// Playlist package implements playlist, -- ordered list of tracks to be played, abstraction.
package playlist

import "container/list"

// Playlist
type Playlist struct {
	// List of tracks present in the current playlist.
	tracks *List
}

// New returns an initialized playlist.
func New() *Playlist {
	return &Playlist{list.New()}
}

// Len returns the total number of tracks present in playlist. 
func (pl *Playlist) Len() int {
	return pl.tracks.Len()
}

// Append insert new track at the end of the playlist.
func (pl *Playlist) Append(track *Track) {
	pl.tracks.PushBack(track)
}

// Clear removes all items from the playlist.
func (pl *Playlist) Clear() {
	pl.tracks.Init()
}
