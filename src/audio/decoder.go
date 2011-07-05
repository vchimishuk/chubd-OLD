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

// decoderFactory is function wich returns new decoder implementation.
type decoderFactory func () Decoder

// List of decoder creation functions.
var decoderFactories []decoderFactory

// RegisterDecoderFactory registers new decoder factory method.
func RegisterDecoderFactory(fact decoderFactory) {
	decoderFactories = append(decoderFactories, fact)
}

// GetDecoder returns decoder for decoding given file.
func GetDecoder(filename string) (decoder Decoder, err os.Error) {
	for _, factory := range decoderFactories {
		decoder = factory()
		if decoder.Match(filename) {
			return decoder, nil
		}
	}

	return nil, os.NewError(fmt.Sprintf("Decoder not found for file '%s'", filename))
}
