package main

import (
	"fmt"
	"net"
)

func startClient() {
	hostName := "example.com" // change this
	portNum := "6000"

	conn, err := net.Dial("tcp", hostName+":"+portNum)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Connection established between %s and localhost.\n", hostName)
	fmt.Printf("Remote Address : %s \n", conn.RemoteAddr().String())
	fmt.Printf("Local Address : %s \n", conn.LocalAddr().String())

}
