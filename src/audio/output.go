// Audio output driver interface.
package audio

import (
	"os"
	"fmt"
)

// Output interface represents audio autput interface (ALSA, OSS, ...).
type Output interface {
	// Open opens output audio device.
	Open() os.Error
	// Set new value for sample rate parameter.
	SetSampleRate(rate int)
	// Set number of channels.
	SetChannels(channels int)
	// Wait waits some free space in output buffer, but not more than maxDelay milliseconds.
	// true result value means that output ready for new portion of data, false -- timeout has occured.
	Wait(maxDelay int) bool
	// AvailUpdate returns free size of output buffer. In bytes.
	AvailUpdate() (size int, err os.Error)
	// Write new portion of data into buffer.
	Write(buf []byte) (written int, err os.Error)
	// Pause pauses playback process.
	Pause()
	// Unpause release pause on playback process.
	Unpause()
	// Close closes output audio device.
	Close()
}

type outputFactory func() Output

// List of avaliable output factories.
var outputFactories map[string]outputFactory = make(map[string]outputFactory)

// RegisterOutput register new output device interface.
func RegisterOutput(name string, fact func() Output) {
	outputFactories[name] = fact
}

// GetOutput returns output interface for writing data to.
func GetOutput() Output {
	// TOOD: Make output configurable. Now alsa hard coded.
	fact, ok := outputFactories["alsa"]
	if !ok {
		panic(fmt.Sprintf("Output '%s' is not available"))
	}

	return fact()
}
