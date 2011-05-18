// Track implements playable file (or piece of file) abstraction.
package filesystem

import (
	"sort"
)

// Track represents track (one song) which can be played.
type Track struct {
	Filename string
}

// NewTrack returns new initialized track.
func NewTrack(filename string) *Track {
	return &Track{filename}
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
	return ta[i].Filename < ta[j].Filename
}

// Swap swaps two array elements.
func (ta TrackArray) Swap(i int, j int) {
	ta[i], ta[j] = ta[j], ta[i]
}

// Sort sorts array in ascending order.
func (ta TrackArray) Sort() {
	sort.Sort(ta)
}
