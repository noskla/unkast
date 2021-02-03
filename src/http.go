package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

// StartHTTPServer launches the main web server for clients/listeners.
func StartHTTPServer(address string, port string) {

	addr := address + ":" + port

	log.Println("Launching web server at", addr)
	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  48 * time.Hour,
		WriteTimeout: 48 * time.Hour,
	}

	http.HandleFunc("/listen", HTTPListenRoute)
	log.Fatalln(httpServer.ListenAndServe())

}

// HTTPListenRoute -> GET /listen?id={uuid.UUID}
// A route which checks if the requested channel in the "id"
// query parameter exist, and if so then sends first 4096
// bytes of the channel's mainBuffer every bitrateInterval
// to synchronize with bitrate.
func HTTPListenRoute(w http.ResponseWriter, r *http.Request) {

	urlQuery, idFound := r.URL.Query()["chan"]
	if !idFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Channel parameter is required."))
		return
	}

	requestedChannel := urlQuery[0]
	if _, ok := ActiveChannels[requestedChannel]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Channel " + requestedChannel + " does not exist."))
		return
	}

	ch := ActiveChannels[requestedChannel]
	ActiveChannels[requestedChannel].status.listeners++
	w.Header().Add("Content-Type", ch.audioFormat)
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	log.Println("Listener started listening to "+requestedChannel+
		". Currently listening:", ActiveChannels[requestedChannel].status.listeners)

	var lastChunk = make([]byte, 4096)
	var currentChunk = make([]byte, 4096)
	var chunkLimitInSecond uint8 = 6
	var chunksPassed uint8 = 0
	var previousSyncTime uint64 = uint64(time.Now().UnixNano())
	var syncSkew int64

	for {
		currentChunk = ch.mainBuffer[0:4096]
		if _, err := w.Write(currentChunk); err != nil ||
			bytes.Compare(lastChunk, currentChunk) == 0 {
			ActiveChannels[requestedChannel].status.listeners--
			break
		}
		lastChunk = currentChunk

		if chunksPassed == chunkLimitInSecond {
			chunksPassed = 0
			timeNow := uint64(time.Now().UnixNano())
			syncSkew = int64(timeNow - (previousSyncTime + uint64(time.Second)))
			previousSyncTime = timeNow
		}
		time.Sleep(ch.bitrateInterval - time.Duration(syncSkew))
		syncSkew = 0
	}

}
