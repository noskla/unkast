package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

// RoutineReadAudio is used in goroutines for each channel
// to read audio files in chunks of 4096 and fill the main
// buffer with data so http server will be able to read from it.
func RoutineReadAudio(filePath string, channelID string) {
	log.Println(channelID+": Reading file", filePath)
	audioFile, err := os.Open(filePath)
	if HandleError(err, false) {
		return
	}
	defer audioFile.Close()

	reader := bufio.NewReader(audioFile)
	moveBuffer := func() {
		ActiveChannels[channelID].mainBuffer = append(
			ActiveChannels[channelID].mainBuffer[4096:],
			ActiveChannels[channelID].temporaryInputBuffer...)
	}

	ActiveChannels[channelID].temporaryInputBuffer = make([]byte, 4096)
	waitInterval := ActiveChannels[channelID].bitrateInterval
	for {
		_, err := reader.Read(ActiveChannels[channelID].temporaryInputBuffer)
		switch err {
		case io.EOF:
			return
		case nil:
			moveBuffer()
		default:
			HandleError(err, false)
			return
		}
		time.Sleep(waitInterval)
	}

}
