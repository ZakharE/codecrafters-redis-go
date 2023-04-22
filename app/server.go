package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()

	request, err := io.ReadAll(conn)
	strings := strings.Split(string(request), " ")
	if len(strings) > 2 {
		writeErr(&conn)
	}

	if strings[0] == "PING" {
		if len(strings) == 1 {
			writeOk(&conn, "PONG")
		} else {
			writeOk(&conn, strings[1])
		}
	}

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}

func writeOk(c *net.Conn, arg string) {
	responseFormat := "+%s\r\n"
	io.WriteString(*c, fmt.Sprintf(responseFormat, arg))
}

func writeErr(conn *net.Conn) {
	io.WriteString(*conn, "-ERR Too many arguments\r\n")
}
