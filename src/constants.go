package unkast

// AudioFormats constant declares the names of
// supported audio formats.
const AudioFormats = map[uint8]string{
	1: "MPEG", 2: "FLAC", 3: "WAV"}

// ChannelTypes constant declares the names of
// supported channel types.
const ChannelTypes = map[uint8]string{
	1: "audioOnly", 2: "streamOnly", 3: "mirror"}
