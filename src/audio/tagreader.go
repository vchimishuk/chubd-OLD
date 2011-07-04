package audio

import (
	"os"
	"fmt"
)

// TagReader interface wraps methods for working with audio file tags.
type TagReader interface {
	// Match returns true it given file can be processed with current TagReader.
	Match(filename string) bool
	// ReadTag parse audio file's metadata and returns filled Tag object.
	ReadTag(filename string) (tag *Tag, err os.Error)
}

// All supported tagreader factory functions.
var readerFactories []func() TagReader

// RegisterTagReaderFactory registers new TagReader factory method.
func RegisterTagReaderFactory(fact func() TagReader) {
	readerFactories = append(readerFactories, fact)
}

// NewTagReader returns TagReader for given file.
func NewTagReader(filename string) (reader TagReader, err os.Error) {
	for _, factory := range readerFactories {
		reader = factory()
		if reader.Match(filename) {
			return reader, nil
		}
	}

	return nil, os.NewError(fmt.Sprintf("TagReadrer not found for file '%s'", filename))
}
