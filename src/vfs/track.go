// Track implements playable file (or piece of file) abstraction.
package vfs

import (
	"sort"
	"./audio"
)

// Track represents track (one song) which can be played.
type Track struct {
	// Full path to the file.
	FilePath *Path
	Number   int
	Tag      *audio.Tag
	// Length of the track in seconds.
	// length int 
}

// NewTrack returns new initialized track indentify some audio file and track.
func NewTrack(filePath *Path, number int) *Track {
	track := new(Track)
	track.FilePath = filePath
	track.Number = number

	return track
}

// Len returns length of the track in seconds.
func (track *Track) Len() int {
	return 0
}

// LenString returns length of the track in standard time format.
func (track *Track) LenString() string {
	return "0:00"
}

// TrackArray is helper type for manipulating Track arrays.
type TrackArray []*Track

// Len returns length of the array.
func (ta TrackArray) Len() int {
	return len(ta)
}

// Less returns true if i-element of the array less than j-element.
func (ta TrackArray) Less(i int, j int) bool {
	return ta[i].FilePath.Path() < ta[j].FilePath.Path()
}

// Swap swaps two array elements.
func (ta TrackArray) Swap(i int, j int) {
	ta[i], ta[j] = ta[j], ta[i]
}

// Sort sorts array in ascending order.
func (ta TrackArray) Sort() {
	sort.Sort(ta)
}
