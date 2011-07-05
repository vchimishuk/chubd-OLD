// playing (decode -> output) thread implementation.
package player

import (
	"os"
	"./vfs"
	"./audio"
)

// messageType is the type for describing messages.
type messageType int

const (
	// Stop current playing and clear tracks queue
	messageTypeStop = iota
	// Queue next track for playing.
	messageTypePlay
	// Toggle pause state (if playing something.
	messageTypePaused
	// Stop goroutine execution request message.
	messageTypeKill
)

// threadState type describes state of the playing thread.
type threadState int

const (
	// Thread is currently stopped.
	threadStateStopped = iota
	// Thread is playing some track now.
	threadStatePlaying
	// Thread is paused.
	threadStatePaused
)

// message type for manipulating playingThread's behaviour.
type message struct {
	t    messageType
	data interface{}
}

// playingThread structure represents thread which decode audio file and writes
// resulting raw PCM data to the output driver. Also is manages some related stuff, like pause, seek, etc.
type playingThread struct {
	// Channel for sending messages to the running goroutine.
	messages chan *message
	// Current state.
	state threadState
	// Separate goroutine will write here values to let us know that we can decode and write new portion of data.
	bufAvailable chan bool
	// Output driver.
	output audio.Output
	// Decoder driver implementation for the current playing track.
	// This is not nil only if thread in threadStatePlaying state.
	decoder audio.Decoder
}

// newPlayingThread returns newly initialized playingThread object.
func newPlayingThread() *playingThread {
	thread := new(playingThread)
	thread.messages = make(chan *message)
	thread.state = threadStateStopped
	thread.bufAvailable = make(chan bool)

	return thread
}

// Start method runs the thread.
func (thread *playingThread) Start() {
	go thread.routine()
}

// Stop release resources and prepare for termination.
func (thread *playingThread) Stop() {
	// We send to the routine message and will wait for the answer.
	wait := make(chan bool)

	msg := new(message)
	msg.t = messageTypeKill
	msg.data = wait
	thread.sendMessage(msg)

	<- wait
}

// Play start playing given track.
func (thread *playingThread) Play(track *vfs.Track) {
	msg := new(message)
	msg.t = messageTypePlay
	msg.data = track
	thread.sendMessage(msg)
}

// Pause toggle pause state.
func (thread *playingThread) Pause() {
	msg := new(message)
	msg.t = messageTypePaused
	thread.sendMessage(msg)
}

// SendMessage queue new message for the playingThread.
func (thread *playingThread) sendMessage(msg *message) {
	thread.messages <- msg
}

// startBufAvailableChecker runs goroutine for checking if output buffer is avaliable
// for new portion of data.
func (thread *playingThread) startBufAvailableChecker() {
	// XXX: I think we can have more than one such goroutine running in time,
	//      But I'll hope it is not a big problem. In other case we can provide mutex
	//      for the nex goroutine duplicate protection.
	if thread.state == threadStatePlaying {
		go func() {
			if thread.output.Wait(500) {
				thread.bufAvailable <- true
			}
		}()
	}
}

// openOutput intializes output driver.
func (thread *playingThread) openOutput() os.Error {
	output := audio.GetOutput()
	err := output.Open()
	if err != nil {
		return err
	}
	output.SetSampleRate(44100)
	output.SetChannels(2)

	thread.output = output

	return nil
}

// closeOutput close and release output driver.
func (thread *playingThread) closeOutput() {
	if thread.output != nil {
		thread.output.Close()
		thread.output = nil
	}
}

// openDecoder initilizes decoder driver.
func (thread *playingThread) openDecoder(track *vfs.Track) os.Error {
	decoder, err := audio.GetDecoder(track.Filename.PathFull())
	if err != nil {
		return err
	}
	thread.decoder = decoder
	thread.decoder.Open(track.Filename.PathFull())

	return nil
}

// closeDecoder releases decoder driver.
func (thread *playingThread) closeDecoder() {
	thread.decoder.Close()
	thread.decoder = nil
}

// Ruotine is the core goroutine function.
func (thread *playingThread) routine() {
	var err os.Error

	for {
		thread.startBufAvailableChecker()

		select {
		case msg := <-thread.messages:
			// Change thread state.
			switch msg.t {
			case messageTypePlay:
				track := msg.data.(*vfs.Track)

				// Initialize decoder driver.
				if thread.decoder != nil {
					thread.closeDecoder()
				}
				err = thread.openDecoder(track)
				if err != nil {
					// TODO: Write into log about unsupported decoder.
					continue // for loop
				}

				// Initialize output driver.
				if thread.output == nil {
					err = thread.openOutput()
					if err != nil {
						// TODO: Write to log.
						continue // for loop
					}
				} else {
					// If this track audio parameters (sample rate, channels number)
					// difference from previous one we need to reconfigure output driver.

					// TODO: Check for the new file format,
					//       maybe we have to change sample rate or channels count.
				}

				thread.state = threadStatePlaying
				thread.startBufAvailableChecker()
			case messageTypePaused:
				if thread.state == threadStatePlaying {
					thread.output.Pause()
					thread.state = threadStatePaused
				} else if thread.state == threadStatePaused {
					thread.output.Unpause()
					thread.state = threadStatePlaying
				}
			case messageTypeStop:
				thread.closeDecoder()
				thread.closeOutput()
				thread.state = threadStateStopped
			case messageTypeKill:
				thread.closeDecoder()
				thread.closeOutput()
				thread.state = threadStateStopped
				return
			}
		case <-thread.bufAvailable:
			// pass
		}

		// Do some job.
		if thread.state == threadStatePlaying {
			size, _ := thread.output.AvailUpdate()
			buf := make([]byte, size)
			thread.decoder.Read(buf)
			thread.output.Write(buf)
		}
	}
}
