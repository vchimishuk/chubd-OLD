// Filesystem package implements abstraction over OS filesystem API.
package vfs

import (
	"os"
	"fmt"
	"path"
	"strings"
	"./config"
)

// Filesystem structure.
type Filesystem struct {
	// Working directory based from the root.
	wd *Path
}

// New returns newly initialized Filesystem object.
func New() *Filesystem {
	fs := new(Filesystem)
	fs.SetWorkingDir("/")

	return fs
}

// WorkingDir returns chrootedd current directory where we are located in.
func (fs *Filesystem) WorkingDir() string {
	return fs.wd.Path()
}

// WorkingDirFull returns not chrooted working directory.
// Returned value is a full path based on root of the physical FS.
func (fs *Filesystem) WorkingDirFull() string {
	return fs.wd.PathFull()
}

// SetWorkingDir sets new working directory, -- directory where we are located in.
func (fs *Filesystem) SetWorkingDir(dir string) os.Error {
	root, _ := config.Configurations.GetString("fs.root")

	var newWd *Path
	if path.IsAbs(dir) {
		newWd = NewPath(dir)
	} else {
		newWd = NewPath(path.Join(fs.wd.Path(), dir))

		// New path can't be upper than root.
		if !strings.HasPrefix(newWd.PathFull(), root) {
			newWd = NewPath("/")
		}
	}

	fileInfo, err := os.Stat(newWd.PathFull())
	if err != nil {
		return os.NewError(fmt.Sprintf("'%s' is not file or directory", newWd.Path()))
	}
	if !fileInfo.IsDirectory() {
		return os.NewError(fmt.Sprintf("'%s' is not a directory", newWd.Path()))
	}

	fs.wd = newWd

	return nil
}

// List returns content of the working directory.
func (fs *Filesystem) List() (entries []*Entry, err os.Error) {
	wd, err := os.Open(fs.WorkingDirFull())
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
		filePath := NewPath(path.Join(fs.WorkingDir(), file))
		fi, err := os.Stat(filePath.PathFull())
		if err != nil {
			return nil, err
		}

		if fi.IsRegular() {
			track, err := NewTrack(filePath)
			if err == nil {
				tracks = append(tracks, track)
			}
		} else if fi.IsDirectory() {
			dir, err := NewDirectory(filePath)
			if err == nil {
				dirs = append(dirs, dir)
			}
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
