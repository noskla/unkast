package main

import "log"

func main() {

	channelID, ok := CreateAudioChannel(192, false, AudioFormats["MPEG"], ChannelTypes["FilesOnly"])
	if !ok {
		log.Fatalln("Error creating channel")
	}

	go RoutineReadAudio("/home/alis/Music/Elektronomia - Magic.mp3", channelID.String())
	StartHTTPServer("127.0.0.1", "8080")

}
