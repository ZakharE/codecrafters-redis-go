package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var CRLFLEN = int32(2)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	//c := make(chan os.Signal, 1)
	////signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	////
	////go func() {
	////	for {
	////		select {
	////		case <-c:
	////			fmt.Println("process was killed")
	////			return
	////		default:
	//
	//		}
	//	}
	//}()

	for {
		conn, err := l.Accept()
		go func() {
			fmt.Println("accepted")
			if err != nil {
				fmt.Println("failed during connection")
			}
			handleConn(conn)
			conn.Close()
		}()
	}
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}

func handleConn(conn net.Conn) {
loop:
	for {
		command, _ := ReadCommand(conn)
		if len(command.Args) > 1 {
			writeErr(&conn)
		}
		switch strings.ToUpper(string(command.Name)) {
		case "PING":
			if len(command.Args) == 1 {
				writeOk(&conn, command.Args[0])
			} else {
				writeOk(&conn, []byte("PONG"))
			}
		case "QUIT":
			break loop
		default:
			writeOk(&conn, []byte("not implemented"))
		}
	}
}

func writeOk(c *net.Conn, arg []byte) {
	responseFormat := "+%s\r\n"
	io.WriteString(*c, fmt.Sprintf(responseFormat, string(arg)))
}

func writeErr(conn *net.Conn) {
	io.WriteString(*conn, "-ERR Too many arguments\r\n")
}

type Command struct {
	Name []byte
	Args [][]byte
}

func ReadCommand(conn net.Conn) (Command, error) {
	buf := make([]byte, 256)

	_, _ = conn.Read(buf)
	println("buf: ", string(buf))

	if buf[0] != '*' {
		return Command{}, errors.New("RESP unsupported format")
	}
	i := int32(1)
	argsNum := int32(0)
	for buf[i] != '\r' {
		argsNum = (argsNum * 10) + int32(buf[i]-'0')
		i++
	}
	array := parseBulkStringArray(buf[i+CRLFLEN:], argsNum)
	return Command{Name: array[0], Args: array[1:]}, nil
}

func parseBulkStringArray(buf []byte, length int32) [][]byte {
	var (
		result [][]byte
	)
	for length > 0 {
		wordLength := int32(0)
		i := int32(1)
		for buf[i] != '\r' {
			wordLength = (wordLength * 10) + int32(buf[i]-'0')
			i++
		}
		i += CRLFLEN
		result = append(result, buf[i:i+wordLength])
		buf = buf[i+wordLength+CRLFLEN:]
		length--
	}
	return result
}
