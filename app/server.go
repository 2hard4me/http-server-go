package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data", err)
		return
	}
	request := string(buf)
	status := strings.Split(request, "\r\n")
	path := strings.Split(status[0], " ")[1]
	method := strings.Split(status[0], " ")[0]
	

	var response []byte
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(path, "/echo") {
		randStr := path[6:]
		response = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len([]rune(randStr))) + "\r\n\r\n" + randStr + "\r\n")
	} else if strings.HasPrefix(path, "/user-agent") {
		userAgent := strings.Split(status[2], " ")[1]
		response = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len([]rune(userAgent))) + "\r\n\r\n" + userAgent + "\r\n")
	} else if strings.HasPrefix(path, "/files") {
		file := strings.Split(path, "/files/")[1]
		dir := os.Args[2]
		
		if method == "GET" {
			data, err := os.ReadFile(filepath.Join(dir, file))
			if err != nil {
				response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
				conn.Write(response)
				return
			}
			response = []byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n\r\n" + string(data) + "\r\n")
		} else if method == "POST" {
			body := strings.Split(request, "\r\n\r\n")[1]
			body = strings.Trim(body, "\x00/")
			err := os.WriteFile(filepath.Join(dir, file), []byte(body), 0644)
			if err != nil {
				fmt.Println("Error writing data to file")
				return
			}
			response = []byte("HTTP/1.1 201 Created\r\n\r\n")
		}
		
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error responding")
		return
	}
}
