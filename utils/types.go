package utils

import (
	"bytes"
	"fmt"
	"strconv"
)

type UploadResponse struct {
	IsError bool
	Message string
	DataId  *string
}

func (r *UploadResponse) ParseResponseToStruct(resBytes []byte) error {

	// Response:
	// true\r\n
	// blablabla\r\n
	// dataId\r\n
	// \r\n
	lines := bytes.Split(resBytes, []byte("\r\n"))

	line := string(lines[0])
	msg := string(lines[1])
	id := string(lines[2])

	isError, _ := strconv.ParseBool(line)

	res := UploadResponse{
		IsError: isError,
		Message: msg,
		DataId:  &id,
	}

	*r = res

	return nil
}

func (r *UploadResponse) ParseResponseToBytes() ([]byte, error) {
	if r.Message == "" {
		return []byte{}, fmt.Errorf("failed to parse empty field")
	}

	s := []byte(fmt.Sprintf("%v\r\n%s\r\n%s\r\n\r\n", r.IsError, r.Message, r.DataId))

	return s, nil
}

type Request struct {
	Operation Operation
	FileData  *[]byte
}

func (r *Request) ParseRequestToStruct(reqBytes []byte) error {
	// Request:
	// SAVE\r\n
	// 01020\r\n
	// \r\n

	lines := bytes.Split(reqBytes, []byte("\r\n"))

	opStr := string(lines[0])
	fileData := lines[1]

	var op Operation = 0

	switch opStr {
	case "0":
		op = SAVE
	case "1":
		op = RETRIEVE
	default:
		return fmt.Errorf("failed pattern matching on operations")
	}

	req := Request{
		Operation: op,
		FileData:  &fileData,
	}

	*r = req

	return nil
}

func (r *Request) ParseRequestToBytes() ([]byte, error) {
	if len(*(r.FileData)) == 0 {
		return []byte{}, fmt.Errorf("failed to parse empty file data")
	}

	s := []byte(fmt.Sprintf("%v\r\n%v\r\n\r\n", r.Operation, r.FileData))

	return s, nil
}

type Operation int

const (
	SAVE Operation = iota
	RETRIEVE
)
