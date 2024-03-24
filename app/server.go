package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

func echoHandle(conn net.Conn, msg string) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n"))
	conn.Write([]byte("Content-Length: " + strconv.Itoa(len(msg)) +"\r\n\r\n"))
	conn.Write([]byte(msg))
	conn.Close()
}

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
			conn.Close()
			break
		}
		buf := make([]byte, 1024)
		conn.Read(buf)
		splitReq := strings.Split(string(buf), " ")
		path := splitReq[1]
		if (path == "/") {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			pathParts := strings.Split(path, "/")
			if (pathParts[1] == "echo") {
				echoHandle(conn, path[6:])
			} else {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
		}
	}
}
