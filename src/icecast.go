package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"
	"time"
)

// SocketClients is a map keyed with the client's
// net.Conn instance that contains the most important
// information about the clients.
var SocketClients = make(map[net.Conn]*SocketClient)

// SocketListen is generally used as a goroutine to listen on TCP
// protocol with given network interface and port and to accept
// any valid connection that are redirected to EmulateIcecast
// function that emulates the Icecast2 protocol.
func SocketListen(address string, port string) {
	listener, err := net.Listen("tcp", address+":"+port)
	HandleError(err, true)
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if HandleError(err, false) {
			log.Println("One of the socket clients has failed to connect", err.Error())
			continue
		}

		clientAddress := connection.RemoteAddr().String()
		log.Println("Client with IP", clientAddress, "has connected via sockets.")
		clientAddressSplit := strings.Split(clientAddress, ":")
		clientIP, clientPort := clientAddressSplit[0], clientAddressSplit[1]

		SocketClients[connection] = &SocketClient{clientIP, clientPort, time.Now(), "", ""}

		go EmulateIcecast(connection)
	}
}

// EmulateIcecast handles most of the socket connection,
// firstly it parses the HTTP headers and collects the most
// important data, then proceeds to listen to data stream
// provided by the client that is being appended into the
// mainBuffer of channel described in PUT request.
func EmulateIcecast(connection net.Conn) {
	var incomingMessage []byte = make([]byte, 0, 4096)
	var temporaryBuffer []byte = make([]byte, 256)

	for !bytes.Contains(incomingMessage, []byte{13, 10, 13, 10}) {
		length, err := connection.Read(temporaryBuffer)
		if err != nil {
			break
		}
		incomingMessage = append(incomingMessage, temporaryBuffer[:length]...)
	}

	if !icecastHandleHeaders(connection, string(incomingMessage)) {
		return
	}

	reader := bufio.NewReader(connection)
	channelID := SocketClients[connection].channel
	var err error

	for {
		if _, err = reader.Read(ActiveChannels[channelID].temporaryInputBuffer); err != nil {
			delete(SocketClients, connection)
			connection.Close()
			return
		}

		ActiveChannels[channelID].mainBuffer = append(ActiveChannels[channelID].mainBuffer[4096:],
			ActiveChannels[channelID].temporaryInputBuffer[:]...)
		time.Sleep(ActiveChannels[channelID].bitrateInterval)
	}
}

// Loops through headers sent by the client, registers all
// required data into the client's struct in map[] and returns
// if EmulateIcecast should continue.
func icecastHandleHeaders(connection net.Conn, headers string) bool {
	for _, header := range strings.Split(headers, "\r\n") {
		headerParts := strings.Split(header, " ")

		if strings.HasPrefix(header, "PUT") {
			SocketClients[connection].channel = headerParts[1][1:] // remove "/" prefix
			if _, ok := ActiveChannels[SocketClients[connection].channel]; !ok {
				connection.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
				delete(SocketClients, connection)
				connection.Close()
				return false
			}
		} else if strings.HasPrefix(header, "Content-Type") {
			if SocketClients[connection].channel == "" {
				connection.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
				delete(SocketClients, connection)
				connection.Close()
				return false
			}
			ActiveChannels[SocketClients[connection].channel].audioFormat = headerParts[1]
		}
	}

	connection.Write([]byte("HTTP/1.1 200 OK\r\n" +
		"Server: " + SocketServerName + "/" + Version + "\r\n" +
		"Connection: Close\r\n" + "Accept-Encoding: identity\r\n" +
		"Allow: GET, SOURCE\r\n" + "Cache-Control: no-cache\r\n" +
		"Pragma: no-cache\r\n" + "Access-Control-Allow-Origin: *\r\n"))

	return true
}
