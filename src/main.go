package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	channelID, _ := CreateAudioChannel(192, false, AudioFormats["MPEG"], ChannelTypes["FilesOnly"])
	//go RoutineReadAudio("/home/alis/Music/Beat Thee.mp3", channelID.String())
	go PlayDirectory("/home/alis/Music/", channelID.String(), true, true)

	channelID, _ = CreateAudioChannel(192, false, AudioFormats["MPEG"], ChannelTypes["StreamOnly"])

	go SocketListen("0.0.0.0", "8081")
	StartHTTPServer("0.0.0.0", "8080")

}
