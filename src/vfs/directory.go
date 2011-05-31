package vfs

import (
	"path"
	"os"
	"fmt"
	"sort"
)

// Directory represents "real" directory in our virtual FS.
type Directory struct {
	// Full path to the directory.
	Filename string
	// Short name, -- last segment.
	Name string
}

// NewDirectory returns newly initialized Directory object.
// filename parameter is full path to this directory.
func NewDirectory(filename string) (dir *Directory, err os.Error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return
	}
	if !fi.IsDirectory() {
		err = os.NewError(fmt.Sprintf("'%s' is not a directory", filename))
		return
	}

	dir = new(Directory)
	dir.Filename = filename
	dir.Name = path.Base(filename)

	return
}

// DirectoryArray is helper type for manipulating (e. g. sorting)  Directory arrays.
type DirectoryArray []*Directory

// Len returns length of the array.
func (da DirectoryArray) Len() int {
	return len(da)
}

// Less returns true if i-element less than j-element.
func (da DirectoryArray) Less(i int, j int) bool {
	return da[i].Name < da[j].Name
}

// Swap swaps two elements.
func (da DirectoryArray) Swap(i int, j int) {
	da[i], da[j] = da[j], da[i]
}

// Sort sorts array in ascending order.
func (da DirectoryArray) Sort() {
	sort.Sort(da)
}
