package unkast

// AudioFormats pseudo-constant declares the names of
// supported audio formats.
var AudioFormats = map[uint8]string{
	1: "MPEG", 2: "FLAC", 3: "WAV"}

// ChannelTypes pseudo-constant declares the names of
// supported channel types.
var ChannelTypes = map[uint8]string{
	1: "audioOnly", 2: "streamOnly", 3: "mirror"}
