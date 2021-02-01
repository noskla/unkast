package main

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// ActiveChannels represents a list of all open audio channels
// that has been created using the CreateAudioChannel function.
// It is being constantly read and written to by channel-assigned
// goroutines to retrieve audio information or write the current
// AudioChannelStatus.
var ActiveChannels = make(map[string]*AudioChannel)

// CreateAudioChannel initializes the channel, appends it to the ActiveChannels
// map and either starts listening for streaming clients or starts broadcasting
// audio files, then returns the UUID of the channel and if the channel creation
// succeded.
func CreateAudioChannel(bitrate uint, hidden bool, audioFormat uint8, channelType uint8) (uuid.UUID, bool) {
	var channel AudioChannel
	var channelStatus AudioChannelStatus = AudioChannelStatus{}

	// https://pkg.go.dev/github.com/google/uuid#NewRandom
	// Documentation recommends using uuid.New(), but it calls the panic
	// if the UUID generation failed. I prefer any error to not be fatal
	// and just return the error if anything happens.
	ID, err := uuid.NewRandom()
	if HandleError(err, false) {
		log.Println("Could not generate UUID on audio channel creation.")
		return [16]byte{}, false
	}

	log.Println("Generated UUID for channel: " + ID.String())

	// Sleep duration in nanoseconds between bursts of 4096 bytes.
	var sleepInterval uint = (uint(time.Second) / (bitrate / 32))

	channel = AudioChannel{
		ID:                   ID,
		temporaryInputBuffer: make([]byte, 4096),
		mainBuffer:           make([]byte, 4096*25),
		audioFormat:          audioFormat,
		bitrate:              bitrate,
		bitrateInterval:      time.Duration(sleepInterval),
		hidden:               hidden,
		channelType:          channelType,
		status:               channelStatus,
	}
	ActiveChannels[ID.String()] = &channel

	// TODO: start goroutine

	return ID, true
}
