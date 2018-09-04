package main

import (
	"encoding/json"
	"log"
	"net"
	"word-tokenize-middleware-socket/core"
	"word-tokenize-middleware-socket/util"
)

func buildUDPConnection() *net.UDPConn {
	udpResolver, _ := net.ResolveUDPAddr("udp", ":5000")
	connection, _ := net.ListenUDP("udp", udpResolver)

	log.Println("UDP server address:", connection.LocalAddr())
	return connection
}

func handleRequest(connection *net.UDPConn) (core.Request, *net.UDPAddr) {
	var buffer [2048]byte
	var request core.Request

	cutPoint, requestAddress, _ := connection.ReadFromUDP(buffer[0:])
	json.Unmarshal(buffer[:cutPoint], &request)
	log.Println("Request:", request)

	return request, requestAddress
}

func buildResponse(connection *net.UDPConn, request core.Request, requestAddress *net.UDPAddr) {
	tokens := util.TextTokenize(request)
	response, _ := json.Marshal(core.Response{Content: tokens})
	connection.WriteToUDP(response, requestAddress)
	log.Println("Response:", tokens)
}

func main() {
	connection := buildUDPConnection()

	// Close the connection when the application closes.
	defer connection.Close()

	for {
		request, requestAddress := handleRequest(connection)
		buildResponse(connection, request, requestAddress)
	}
}