// Filesystem package implements abstraction over OS filesystem API.
package vfs

import (
	"os"
	"fmt"
	"path"
	"strings"
	"strconv"
	"cue"
	"./audio"
	"./config"
)

const (
	PlaylistName = "*vfs*"
)

const (
	CueFilesExtension = ".cue"
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

// getDirs returns sorted list of directories in the working folder.
func (fs *Filesystem) getDirs() (dirs []*Directory, err os.Error) {
	wd, err := os.Open(fs.WorkingDirFull())
	if err != nil {
		return nil, err
	}
	defer wd.Close()

	dirnames, err := wd.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	dirs = make([]*Directory, 0, len(dirnames))
	for _, file := range dirnames {
		filePath := NewPath(path.Join(fs.WorkingDir(), file))
		fi, err := os.Stat(filePath.PathFull())
		if err != nil {
			return nil, err
		}

		if fi.IsDirectory() {
			dir, err := NewDirectory(filePath)
			if err == nil {
				dirs = append(dirs, dir)
			} else {
				// TODO: Write waring into log.
			}
		}
	}

	DirectoryArray(dirs).Sort()

	return dirs, nil
}

// newTrack returns new track structure based on the given file.
func (fs *Filesystem) newTrack(path *Path) (track *Track, err os.Error) {
	tagReader, err := audio.NewTagReader(path.PathFull())
	if err != nil {
		return nil, err
	}

	tag, err := tagReader.ReadTag(path.PathFull())
	if err != nil {
		return nil, err
	}

	track = NewTrack(path, 0)
	track.Tag = tag

	return
}

// getTracks returns list of the tracks in the current folder.
func (fs *Filesystem) getTracks() (tracks []*Track, err os.Error) {
	wd, err := os.Open(fs.WorkingDirFull())
	if err != nil {
		return nil, err
	}
	defer wd.Close()

	dirnames, err := wd.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	// audioFiles is the list of regular audio files, they have no any relation to cue.
	audioFiles := make([]*Path, 0, len(dirnames))
	cueFiles := make([]*Path, 0, len(dirnames))
	for _, file := range dirnames {
		filePath := NewPath(path.Join(fs.WorkingDir(), file))
		fi, err := os.Stat(filePath.PathFull())
		if err != nil {
			return nil, err
		}

		if fi.IsRegular() {
			ext := strings.ToLower(path.Ext(filePath.Path()))

			if ext == CueFilesExtension {
				cueFiles = append(cueFiles, filePath)
			} else {
				// If this is supported audio file.
				_, err := audio.NewTagReader(filePath.PathFull())
				if err == nil {
					audioFiles = append(audioFiles, filePath)
				}
			}
		}
	}

	PathArray(cueFiles).Sort()
	PathArray(audioFiles).Sort()

	tracks = make([]*Track, 0, 10)

	// Process cue files firs.
	for _, cueFile := range cueFiles {
		file, err := os.Open(cueFile.PathFull())
		if err != nil {
			panic(err) // TOOD: Write log error message.
			continue
		}

		cueSheet, err := cue.Parse(file)
		if err != nil {
			panic(err) // TODO: Write log error message.
			continue
		}

		for _, cueFile := range cueSheet.Files {
			// Remove this file from the audioFiles
			fileName := path.Base(cueFile.Name)
			for i := 0; i < len(audioFiles); i++ {
				audioFileName := path.Base(audioFiles[i].Path())
				if audioFileName == fileName {
					audioFiles[i] = nil
					break
				}
			}

			filePath := NewPath(path.Join(fs.WorkingDir(), cueFile.Name))

			// Check if we can decode this file.
			_, err := audio.NewTagReader(filePath.PathFull())
			if err != nil {
				panic(err) // TODO: Write warning message to log.
			}

			for _, cueTrack := range cueFile.Tracks {
				track := NewTrack(filePath, cueTrack.Number)
				track.Tag = new(audio.Tag)
				if len(cueTrack.Performer) > 0 {
					track.Tag.Artist = cueTrack.Performer
				} else {
					track.Tag.Artist = cueSheet.Performer
				}
				track.Tag.Album = cueSheet.Title
				track.Tag.Title = cueTrack.Title
				track.Tag.Number = strconv.Itoa(cueTrack.Number)

				tracks = append(tracks, track)
			}
		}
	}

	newAudioFiles := make([]*Path, 0, len(audioFiles))
	for _, audioFile := range audioFiles {
		if audioFile != nil {
			newAudioFiles = append(newAudioFiles, audioFile)
		}
	}
	audioFiles = newAudioFiles

	// Process non cue audio files.
	for _, audioFile := range audioFiles {
		track, err := fs.newTrack(audioFile)
		if err != nil {
			panic(err) // TODO: Write log warning.
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// List returns content of the working directory.
func (fs *Filesystem) List() (entries []*Entry, err os.Error) {
	dirs, err := fs.getDirs()
	if err != nil {
		return nil, fmt.Errorf("Directory listing failed. %s", err.String())
	}

	tracks, err := fs.getTracks()
	if err != nil {
		return nil, fmt.Errorf("Directory listing failed. %s", err.String())
	}

	entries = make([]*Entry, 0, len(dirs)+len(tracks))
	for _, dir := range dirs {
		entries = append(entries, NewEntry(TypeDirectory, dir))
	}
	for _, track := range tracks {
		entries = append(entries, NewEntry(TypeTrack, track))
	}

	return entries, nil
}
