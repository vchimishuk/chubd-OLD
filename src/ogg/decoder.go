package ogg

import (
	"os"
	"fmt"
	ogggo "ogg"
	"./audio"
	"./utils"
)

// Ogg decoder implementation.
type Decoder struct {
	oggFile *ogggo.File
}

func NewDecoder() audio.Decoder {
	decoder := new(Decoder)

	return decoder
}

// See audio.Decoder.
func (decoder *Decoder) Match(filename string) bool {
	return utils.ExtensionMatch(filename, Extension)
}

// See audio.Decoder.
func (decoder *Decoder) Open(filename string) os.Error {
	file, err := ogggo.New(filename)
	if err != nil {
		return os.NewError(fmt.Sprintf("Failed to open ogg decoder. %s", err))
	}

	decoder.oggFile = file

	return nil
}

// See audio.Decoder.
func (decoder *Decoder) Read(buf []byte) (read int, err os.Error) {
	read = decoder.oggFile.Read(buf)

	return read, nil
}

// See audio.Decoder.
func (decoder *Decoder) Close() {
	decoder.oggFile.Close()
}
