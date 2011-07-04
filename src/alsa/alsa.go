// alsa output driver implementation.
package alsa

import (
	alsago "alsa"
	"os"
	"fmt"
	"./audio"
)

// DriverName is the string name of the alsa driver.
var DriverName string = "alsa"

// Alsa aoutput driter handler structure.
type Alsa struct {
	handle *alsago.Handle
}

// New returns newly initialized alsa output driver.
func New() audio.Output {
	return new(Alsa)
}

func (a *Alsa) Open() os.Error {
	a.handle = alsago.New()
	err := a.handle.Open("default", alsago.StreamTypePlayback, alsago.ModeBlock)
	if err != nil {
		return os.NewError(fmt.Sprintf("Failed to open audio output device. %s", err))
	}

	a.handle.SampleFormat = alsago.SampleFormatS16LE // XXX:

	return nil
}

func (a *Alsa) SetSampleRate(rate int) {
	a.handle.SampleRate = rate
	a.handle.ApplyHwParams()
}

func (a *Alsa) SetChannels(channels int) {
	a.handle.Channels = channels
	a.handle.ApplyHwParams()
}

func (a *Alsa) Wait(maxDelay int) bool {
	ok, _ := a.handle.Wait(maxDelay)

	return ok
}

func (a *Alsa) AvailUpdate() (size int, err os.Error) {
	return a.handle.AvailUpdate()
}

func (a *Alsa) Write(buf []byte) (written int, err os.Error) {
	return a.handle.Write(buf)
}

func (a *Alsa) Pause() {
	a.handle.Pause()
}

func (a *Alsa) Unpause() {
	a.handle.Unpause()
}

func (a *Alsa) Close() {
	a.handle.Close()
}
