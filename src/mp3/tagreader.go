package mp3

import (
	"os"
	"path"
	"strings"
	"id3tag"
	"./audio"
)

type TagReader struct {

}

func NewTagReader() *TagReader {
	return new(TagReader)
}

func (tr *TagReader) Match(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))

	return ext == ".mp3"
}

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
