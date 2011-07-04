// Tag reader implementation for ogg files support.
package ogg

import (
	"os"
	"strings"
	"ogg"
	"./audio"
)

// Ogg TagReader implementation.
type TagReader struct {

}

// NewTagreader returns newly initialized ogg TagReader implementation.
func NewTagReader() audio.TagReader {
	return new(TagReader)
}

// Match returns true if given file is the supported ogg file.
func (tr *TagReader) Match(filename string) bool {
	return match(filename)
}

// ReadTag returns Tag structure filled with values from the given file.
func (tr *TagReader) ReadTag(filename string) (tag *audio.Tag, err os.Error) {
	/*
	 User comments example:
	 TITLE=Baby, Please Don't Go
	 ARTIST=AC/DC
	 ALBUM=Some album
	 DATE=1988
	 TRACKNUMBER=01
	 GENRE=Hard Rock
	 DESCRIPTION=some comment here...
	 COMMENT=some comment here...
	 =some comment here...
	*/

	file, err := ogg.New(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tag = new(audio.Tag)

	for _, uc := range file.Comment().UserComments {
		foo := strings.Split(uc, "=", 2)
		key := foo[0]
		value := foo[1]

		switch key {
		case "ARTIST":
			tag.Artist = value
		case "ALBUM":
			tag.Album = value
		case "TITLE":
			tag.Title = value
		}
	}

	tag.Length = "0:00"

	return tag, nil
}
