package main

import (
	"fmt"
	. "go-file-streamer/utils"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":2174")

	if err != nil {
		panic("failed to create tcp conn to server")
	}

	fileData := GetFile()

	r := &Request{
		Operation: 0,
		FileData:  fileData,
	}

	reqByte, err := r.ParseRequestToBytes()

	if err != nil {
		panic(err.Error())
	}

	_, err = conn.Write(reqByte)
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}

	if err != nil {
		panic("failed to write data to tcp conn")
	}

	res := ConnectionLoop(conn)

	fmt.Println("Result:")
	fmt.Printf("Error: %v\n", res.IsError)
	fmt.Printf("Message: %v\n", res.Message)

}

func ConnectionLoop(conn net.Conn) UploadResponse {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	data, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("failed to read response:", err)
		return UploadResponse{}
	}

	response := UploadResponse{}
	err = response.ParseResponseToStruct(data)
	if err != nil {
		panic(err.Error())
	}

	return response
}

func GetFile() *[]byte {
	fileReader, err := os.Open("./assets/pexel_img.jpg")

	if err != nil {
		panic("failed to read image data")
	}

	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)

	if err != nil {
		panic("failed to process image data")
	}

	return &data
}
