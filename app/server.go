package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

type httpRequestParts struct {
	method string
	target string
	hostname string
	useragent string
}

func rootHandle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Close()
}

func echoHandle(conn net.Conn, msg string) {
	fmt.Println("processing echo :" + msg )
	fmt.Println("done processing")
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n"))
	conn.Write([]byte("Content-Length: " + strconv.Itoa(len(msg)) +"\r\n\r\n"))
	conn.Write([]byte(msg))
	conn.Close()
}

func userAgentHandle(conn net.Conn, useragent string) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n"))
	conn.Write([]byte("Content-Length: " + strconv.Itoa(len(useragent)) +"\r\n\r\n"))
	conn.Write([]byte(useragent))
	conn.Close()
}

func defaultHandle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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
		var hrp httpRequestParts
		conn.Read(buf)
		splitLines := strings.Split(string(buf), "\r\n")
		splitLine1 := strings.Split(splitLines[0], " ")
		splitLine2 := strings.Split(splitLines[1], " ")
		splitLine3 := strings.Split(splitLines[2], " ")
		hrp.method = splitLine1[0]
		hrp.target = splitLine1[1]
		hrp.hostname = splitLine2[1]
		hrp.useragent = splitLine3[1]
		pathType := strings.Split(hrp.target, "/")[1]
	//  fmt.Println("method:" + hrp.method)
	//	fmt.Println("target:" + hrp.target)
	//	fmt.Println("hostname:" + hrp.hostname)
	//	fmt.Println("useragent:" + hrp.useragent)
	//	fmt.Println(pathType)

		switch pathType {
		case "":
			go rootHandle(conn)
		case "echo":
			go echoHandle(conn, hrp.target[6:])
		case "user-agent":
			go userAgentHandle(conn, hrp.useragent)
		default:
			go defaultHandle(conn)
		}
	}
}
