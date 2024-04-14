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
	randomStr := path[5:]
	strLen := strconv.Itoa(len([]rune(randomStr)))
	response := "HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length: " + strLen + "\r\n\r\n" + randomStr + "\r\n"

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error responding")
		return
	}
}
