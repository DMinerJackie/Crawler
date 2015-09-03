package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(c net.Conn) {

}

func startServer() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server up and listening on port 6000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
