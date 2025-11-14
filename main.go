package main

import (
	"bufio"
	"fmt"
	"go-file-streamer/handlers"
	"go-file-streamer/storage"
	. "go-file-streamer/utils"
	"io"
	"net"
)

const (
	MAX_FILE_DATA uint = 30_000_000
)

type Config struct {
	lis   net.Listener
	saver storage.Saver
}

var conf Config

func init() {
	ser := storage.NewFileSaver("./db")

	saver := storage.NewSaverService(ser)
	conf.saver = saver

	err := conf.saver.BuildUpCache()

	if err != nil {
		panic(err)
	}
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

	if req.Operation == SAVE {
		fmt.Printf("FileData Length: %v\n", len(req.FileData))
	} else if req.Operation == RETRIEVE {
		fmt.Printf("FileId: %s\n", req.FileId)
	}

	h := handlers.NewHandler()

	switch req.Operation {
	case SAVE:
		b, err := h.Upload(&req.FileData, conf.saver)

		if err != nil {
			h.ErrorResponseWriter(conn, "failed to upload file")
			return
		}

		h.ResponseWriter(conn, b)

		fmt.Println("\nSent response...")
		return

	case RETRIEVE:
		if req.FileId == "" {
			h.ErrorResponseWriter(conn, "failed to find file with missing id")
			return
		}

		b, err := h.Download(req.FileId, conf.saver)

		if err != nil {
			h.ErrorResponseWriter(conn, "failed to download file")
			return
		}

		h.ResponseWriter(conn, b)
		return
	default:
		h.ErrorResponseWriter(conn, "failed pattern matching on operation")
		return
	}

}
