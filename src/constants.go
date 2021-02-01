package main

// AudioFormats pseudo-constant declares the names of
// supported audio formats.
var AudioFormats = map[string]uint8{
	"MPEG": 1, "MP3": 1, "OGG": 2, "AAC": 3}

// ChannelTypes pseudo-constant declares the names of
// supported channel types.
var ChannelTypes = map[string]uint8{
	"FilesOnly": 1, "StreamOnly": 2, "Mirror": 3}
