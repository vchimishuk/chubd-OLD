package ogg

import (
	"strings"
	"path"
)

// match returns true is given file is the valid ogg file.
func match(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))

	return ext == ".ogg"
}
