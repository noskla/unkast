package unkast

import (
	"log"
	"net/http"
	"time"
)

// StartHTTPServer launches the main web server for clients/listeners.
func StartHTTPServer(address string, port string) {

	addr := address + ":" + port

	log.Println("Launching web server at ", addr)
	httpServer := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       48 * time.Hour,
		WriteTimeout:      5 * time.Second,
	}

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
	for {
		if _, err := w.Write(ch.mainBuffer[:4096]); err != nil {
			break
		}
		time.Sleep(ch.bitrateInterval)
	}

}
