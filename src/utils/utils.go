package utils

import (
	"strings"
	"path"
)

// ExtensionMatch returns true is given file's extension matches with pattern.
// Notice, that extension should be in the lower case. E. g. ".ogg"
func ExtensionMatch(filename string, extension string) bool {
	ext := strings.ToLower(path.Ext(filename))

	return ext == extension
}
