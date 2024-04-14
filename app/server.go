package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
	} else if strings.HasPrefix(path, "/echo") {
		randStr := path[5:]
		response = []byte("HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len([]rune(randStr))) + "\r\n\r\n" + randStr + "\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error responding")
		return
	}
}
