package mp3

import (
	"os"
	"./audio"
	"./utils"
)

// MP3 decoder implementation.
type Decoder struct {

}

// NewDecoder returns MP3 decoder implementation.
func NewDecoder() audio.Decoder {
	return new(Decoder)
}

// See audio.Decoder.
func (decoder *Decoder) Match(filename string) bool {
	return utils.ExtensionMatch(filename, Extension)
}

// See audio.Decoder.
func (decoder *Decoder) Open(filename string) os.Error {
	// TODO: Not implemeted.
	return os.NewError("Not implemented")
}

// See audio.Decoder.
func (decoder *Decoder) Read(buf []byte) (read int, err os.Error) {
	// TODO: Not implemeted.
	return -1, os.NewError("Not implemented")
}

// See audio.Decoder.
func (decoder *Decoder) Close() {
	// TODO: Not implemeted.
}
