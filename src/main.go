package main

import (
	"math/rand"
	"time"
	"os"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	channelID, _ := CreateAudioChannel(192, false, AudioFormats["MPEG"], ChannelTypes["FilesOnly"])
	//go RoutineReadAudio("/home/alis/Music/Beat Thee.mp3", channelID.String())
	homedir, _ := os.UserHomeDir()
	var playdir = homedir + "/Music/"
	if len(os.Args) > 1 {
		playdir = os.Args[1]
	}
	go PlayDirectory(playdir, channelID.String(), true, true)

	channelID, _ = CreateAudioChannel(192, false, AudioFormats["MPEG"], ChannelTypes["StreamOnly"])

	go SocketListen("0.0.0.0", "8081")
	StartHTTPServer("0.0.0.0", "8080")

}
