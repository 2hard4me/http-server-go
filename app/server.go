package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data")
		return
	}
	request := string(buf)
	status := strings.Split(request, "\r\n")
	path := strings.Split(status[0], " ")[1]

	var response []byte
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error responding")
		return
	}
}
