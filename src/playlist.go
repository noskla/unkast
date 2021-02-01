package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
)

// PlayDirectory is generally ran as a goroutine, it lists the
// directory and passes the filenames into the RoutineReadAudio
// function.
func PlayDirectory(directoryPath string, channelID string, loop bool, shuffle bool) {
	if !strings.HasPrefix(directoryPath, string(os.PathSeparator)) {
		directoryPath += string(os.PathSeparator)
	}

	directory, err := os.Open(directoryPath)
	if HandleError(err, false) {
		log.Println("Failure opening directory.")
		return
	}

	files, err := directory.Readdirnames(math.MaxUint16)
	if HandleError(err, false) {
		log.Println("Failure listing directory.")
		return
	}

	if shuffle {
		rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	}

	for _, audioFile := range files {
		RoutineReadAudio(directoryPath+audioFile, channelID)
	}

	if loop {
		PlayDirectory(directoryPath, channelID, loop, shuffle)
	}
}
