package vfs

import (
	. "path"
	"strings"
	"./config"
)
// Path represents path to the resource of the VFS (file, track, ...)
// and helps to convert OS FS pathes into root-based VFS ones.
type Path struct {
	// path is the VFS root-based path.
	rootedPath string
}

// NewPath returns newly initialized Path object
// for given VFS (based on fs.root) path value.
func NewPath(filename string) *Path {
	return &Path{filename}
}

// NewPathFull returns newly initialized Path object
// for given OS-like path value.
func NewPathFull(filename string) *Path {
	root, _ := config.Configurations.GetString("fs.root")

	if !strings.HasSuffix(filename, root) {
		panic("NewPathRooted should be used insted of NewPath")
	}

	return NewPathFull(filename[len(root):])
}

// String returns string representation of the object.
func (path *Path) String() string {
	return path.Path()
}

// Path returns VFS (based on fs.root) path value.
func (path *Path) Path() string {
	return path.rootedPath
}

// PathFull return full physical path (as used in OS).
func (path *Path) PathFull() string {
	root, _ := config.Configurations.GetString("fs.root")

	return Join(root, path.rootedPath)
}
