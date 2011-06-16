// Tag reader implementation for mp3 files support.
package mp3

import (
	"os"
	"path"
	"strings"
	"id3tag"
	"./audio"
)

// MP3 TagReader implementation.
type TagReader struct {

}

// NewTagreader returns newly initialized MP3 TagReader implementation.
func NewTagReader() *TagReader {
	return new(TagReader)
}

// Match returns true if given file is the supported MP3 file.
func (tr *TagReader) Match(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))

	return ext == ".mp3"
}

// ReadTag returns Tag structure filled with values from the given MP3 file.
func (tr *TagReader) ReadTag(filename string) (tag *audio.Tag, err os.Error) {
	id3Tag, err := id3tag.Parse(filename)
	if err != nil {
		return nil, err
	}

	tag = new(audio.Tag)
	tag.Artist = id3Tag.Artist()
	tag.Album = id3Tag.Album()
	tag.Title = id3Tag.Title()
	tag.Length = "0:00"

	return tag, nil
}
