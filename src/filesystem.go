// Filesystem package implements abstraction over OS filesystem API.
package filesystem

import (
	"os"
	"fmt"
	"path"
)

// Filesystem structure.
type Filesystem struct {
	// Working directory.
	wd string
}

// New returns newly initialized Filesystem object.
func New() *Filesystem {
	fs := new(Filesystem)
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize new Filesystem: %s", err))
	}

	fs.SetWorkingDir(wd)

	return fs
}

// WorkingDir returns current directory where we are located in.
func (fs *Filesystem) WorkingDir() string {
	return fs.wd
}

// SetWorkingDir sets new working directory, -- directory where we are located in.
func (fs *Filesystem) SetWorkingDir(dir string) os.Error {
	var newWd string
	dir = path.Clean(dir)

	if path.IsAbs(dir) {
		newWd = dir
	} else {
		newWd = path.Join(fs.wd, dir)
	}

	fileInfo, err := os.Stat(newWd)
	if err != nil {
		return err
	}

	if !fileInfo.IsDirectory() {
		return os.NewError(fmt.Sprintf("'%s' is not a directory", newWd))
	}

	fs.wd = newWd

	// TODO: Check if directory is readable end executable for me.

	return nil
}

// List returns content of the working directory.
func (fs *Filesystem) List() (entries []*Entry, err os.Error) {
	wd, err := os.Open(fs.wd)
	if err != nil {
		return nil, err
	}
	defer wd.Close()

	files, err := wd.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	dirs := make([]*Directory, 0, len(files))
	tracks := make([]*Track, 0, len(files))

	// Split items to dirs and tracks.
	for _, file := range files {
		fullPath := path.Join(fs.wd, file)
		fi, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}

		if fi.IsRegular() {
			tracks = append(tracks, NewTrack(fullPath))
		} else if fi.IsDirectory() {
			d, _ := NewDirectory(fullPath)
			dirs = append(dirs, d)
		}
	}

	// Sort this two parts
	DirectoryArray(dirs).Sort()
	TrackArray(tracks).Sort()

	// End join them all.
	entries = make([]*Entry, 0, len(dirs)+len(tracks))
	for _, dir := range dirs {
		entries = append(entries, NewEntry(TypeDirectory, dir))
	}
	for _, track := range tracks {
		entries = append(entries, NewEntry(TypeTrack, track))
	}

	return entries, nil
}
