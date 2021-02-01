package unkast

import (
	"time"

	"github.com/google/uuid"
)

// AudioChannel represents a unique channel with it's own
// goroutine in which audio files are being read into
// the temporaryInputBuffer in 4096 byte chunks and then appended
// into the mainBuffer, which is generally initialized with size
// of (25 * 4096). If channelType is {ChannelTypes[2] ("streamOnly")}
// then the channel will not read files and will wait for incoming
// streams from supported streaming clients.
type AudioChannel struct {
	ID                   uuid.UUID
	temporaryInputBuffer []byte
	mainBuffer           []byte
	audioFormat          uint8
	bitrate              uint
	bitrateInterval      time.Duration
	hidden               bool
	channelType          uint8
	status               AudioChannelStatus
}

// AudioChannelStatus is a constantly changing part of AudioChannel,
// indicating the current state of the channel, including streamName
// - a general name of the broadcast, it's description, listener count
// and currentState, which can be used to pass a currently played song
// title and it's artist.
type AudioChannelStatus struct {
	streamName        string
	streamDescription string
	listeners         uint
	currentState      string
}
