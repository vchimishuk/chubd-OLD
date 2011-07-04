// Audio decoder interface.
package audio

import (
	"os"
	"fmt"
)

// Decoder interface represents audio decoder for the particular audio format.
type Decoder interface {
	// Match returns true if given file supported by this decoder.
	Match(filename string) bool
	// Open inialize decoder object.
	Open(filename string) os.Error
	// Read decode piece of data and returns raw PCM audio data.
	Read(buf []byte) (read int, err os.Error)
	// Close releases decoder resources.
	Close()
}

// List of decoder creation functions.
var decoderFactories []func() Decoder

// RegisterDecoderFactory registers new decoder factory method.
func RegisterDecoderFactory(fact func() Decoder) {
	decoderFactories = append(decoderFactories, fact)
}

func NewDecoder(filename string) (decoder Decoder, err os.Error) {
	for _, factory := range decoderFactories {
		decoder = factory()
		if decoder.Match(filename) {
			return decoder, nil
		}
	}

	return nil, os.NewError(fmt.Sprintf("Decoder not found for file '%s'", filename))
}
