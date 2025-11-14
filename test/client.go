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
	testUpload()
	testDownload()
}

func testUpload() {
	conn, err := net.Dial("tcp", ":2174")

	if err != nil {
		panic("failed to create tcp conn to server")
	}

	fileData := GetFile()

	// Request File Retrieve
	r := &Request{
		Operation: 0,
		FileData:  *fileData,
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

	resRetrieve := ConnectionLoopUpload(conn)

	fmt.Println("Result:")
	fmt.Printf("Error: %v\n", resRetrieve.IsError)
	fmt.Printf("Message: %v\n", resRetrieve.Message)
	fmt.Printf("File Data: %v\n", resRetrieve.DataId)

	fmt.Println("\nUpload finished...")
}

func testDownload() {
	conn, err := net.Dial("tcp", ":2174")

	if err != nil {
		panic("failed to create tcp conn to server")
	}

	// conn, _ = net.Dial("tcp", ":2174")

	// Request File Retrieve
	tmp := "ff44ebc5-b8fb-4194-834c-5731907a9f55"
	r := &Request{
		Operation: 1,
		FileId:    tmp,
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

	resRetrieve := ConnectionLoopRetrieve(conn)

	fmt.Println("Result:")
	fmt.Printf("Error: %v\n", resRetrieve.IsError)
	fmt.Printf("Message: %v\n", resRetrieve.Message)
	fmt.Printf("File Length: %v\n", len(*resRetrieve.FileData))

	fmt.Println("\nDownload finished...")

}

func ConnectionLoopUpload(conn net.Conn) UploadResponse {
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

func ConnectionLoopRetrieve(conn net.Conn) RetrieveResponse {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	data, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("failed to read response:", err)
		return RetrieveResponse{}
	}

	response := RetrieveResponse{}
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
