package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	method  string
	path    string
	version string
	headers map[string]string
}

func main() {	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading buffer: ", err.Error())
	}

	request := parseRequest(string(buffer))

	if request.path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if (strings.Contains(request.path, "/echo/")) {
		str := request.path[len("/echo/"):]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
		conn.Write([]byte(response))
	} else if (strings.Contains(request.path, "/user-agent")) {
		str := request.headers["user-agent"]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
		conn.Write([]byte(response))
		} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func parseRequest(requestString string) HttpRequest {
	request := HttpRequest{
		headers: make(map[string]string),
	}

	lines := strings.Split(requestString, "\r\n")
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) >= 3 {
		request.method = requestLine[0]
		request.path = requestLine[1]
		request.version = requestLine[2]
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}

		if colonIndex := strings.Index(line, ": "); colonIndex != -1 {
			key := strings.ToLower(line[:colonIndex])
			value := line[colonIndex + len(": "):]
			request.headers[key] = value
		}
	}

	return request
}