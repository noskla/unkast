package unkast

import (
	"bufio"
	"io"
	"os"
)

// RoutineReadAudio is used in goroutines for each channel
// to read audio files in chunks of 4096 and fill the main
// buffer with data so http server will be able to read from it.
func RoutineReadAudio(filePath string, channelID string) {
	audioFile, err := os.Open(filePath)
	if HandleError(err, false) {
		return
	}
	defer audioFile.Close()

	reader := bufio.NewReader(audioFile)
	moveBuffer := func(chunk []byte) {
		ActiveChannels[channelID].mainBuffer = append(
			ActiveChannels[channelID].mainBuffer[4096:], chunk...)
	}

	for {
		_, err = reader.Read(ActiveChannels[channelID].temporaryInputBuffer)
		switch err {
		case io.EOF:
			break
		case nil:
			moveBuffer(ActiveChannels[channelID].temporaryInputBuffer)
		default:
			HandleError(err, false)
			return
		}
	}

}
