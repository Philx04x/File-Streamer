package main

import (
	"bufio"
	"fmt"
	"go-file-streamer/handlers"
	. "go-file-streamer/utils"
	"io"
	"net"
)

const (
	MAX_FILE_DATA uint = 30_000_000
)

type Config struct {
	lis net.Listener
}

func main() {
	server := NewTCPServer()

	chn := make(chan bool)
	go AcceptLoop(server)
	<-chn
}

func NewTCPServer() net.Listener {
	lis, err := net.Listen("tcp", ":2174")

	if err != nil {
		panic("failed to serve tcp connection")
	}

	return lis
}

func AcceptLoop(lis net.Listener) {
	for {
		conn, err := lis.Accept()

		if err != nil {
			fmt.Println("failed to connect")

			conn.Close()
			return
		}

		go HandleLoop(conn)
	}
}

func HandleLoop(conn net.Conn) {
	bufio := bufio.NewReader(conn)
	buffer := make([]byte, 0, 1024)

	var idx uint = 0

	for {
		dataByte, err := bufio.ReadByte()

		if idx == MAX_FILE_DATA {
			fmt.Println("failed to read byte. Buffer overflow")
			return
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("failed to parse line")

			conn.Close()
			return
		}

		buffer = append(buffer, dataByte)

		idx++
	}

	req := Request{}
	err := req.ParseRequestToStruct(buffer)

	if err != nil {
		fmt.Println(err.Error())

		conn.Close()
		return
	}

	fmt.Println("Request Data:")
	fmt.Printf("Operation: %v\n", req.Operation)
	fmt.Printf("FileData Length: %v\n", len(*(req.FileData)))

	h := handlers.NewHandler()

	switch req.Operation {
	case SAVE:
		b, err := h.Upload(req.FileData)

		if err != nil {
			h.ErrorResponseWriter(conn, "failed to upload file")
			return
		}

		h.ResponseWriter(conn, b)

		fmt.Println("\nSent response...")
		return

	case RETRIEVE:
		h.ErrorResponseWriter(conn, "operation not implemented")
		return
	default:
		h.ErrorResponseWriter(conn, "failed pattern matching on operation")
		return
	}

}
