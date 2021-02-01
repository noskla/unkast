package main

// AudioFormats pseudo-constant declares the names of
// supported audio formats.
var AudioFormats = map[string]string{
	"MPEG": "audio/mpeg", "MP3": "audio/mpeg", "OGG": "audio/ogg", "AAC": "audio/aac"}

// ChannelTypes pseudo-constant declares the names of
// supported channel types.
var ChannelTypes = map[string]uint8{
	"FilesOnly": 1, "StreamOnly": 2, "Mirror": 3}
