package audio

import (
	"os"
	"fmt"
)

// All supported tagreaders.
var readers = make([]TagReader, 0, 0)

// TagReader interface wraps methods for working with audio file tags.
type TagReader interface {
	// Match returns true it given file can be processed with current TagReader.
	Match(filename string) bool
	// ReadTag parse audio file's metadata and returns filled Tag object.
	ReadTag(filename string) (tag *Tag, err os.Error)
}

// RegisterTagReader registers new TagReader interface implementation.
// For example, before you can read ID3 tags you need to register reader
// that supports ID3 tags.
func RegisterTagReader(reader TagReader) {
	readers = append(readers, reader)
}

// NewTagReader returns TagReader for given file.
func NewTagReader(filename string) (reader TagReader, err os.Error) {
	for i := 0; i < len(readers); i++ {
		if readers[i].Match(filename) {
			return readers[i], nil
		}
	}

	return nil, os.NewError(fmt.Sprintf("TagReadrer not found for file '%s'", filename))
}
