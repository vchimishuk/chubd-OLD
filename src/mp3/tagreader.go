package mp3

import (
	"os"
	"path"
	"strings"
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
	tag = new(audio.Tag)
	tag.Artist = "artist"
	tag.Album = "album"
	tag.Title = "title"
	tag.Length = "0:00"

	return tag, nil
}
