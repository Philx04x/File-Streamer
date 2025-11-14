package utils

import (
	"bytes"
	"fmt"
)

func (r *Request) ParseRequestToStruct(reqBytes []byte) error {
	// Request:
	// SAVE\r\n
	// 01020\r\n
	// \r\n

	lines := bytes.Split(reqBytes, []byte("\r\n"))

	opStr := string(lines[0])

	var op Operation = 0

	switch opStr {
	case "0":
		op = SAVE
	case "1":
		op = RETRIEVE
	default:
		return fmt.Errorf("failed pattern matching on operations")
	}

	var fileData []byte = []byte{}
	var fileId string = ""

	if op == SAVE {
		fileData = lines[1]
	} else if op == RETRIEVE {
		fileId = string(lines[1])
	} else {
		return fmt.Errorf("failed if matching on operations")
	}

	req := Request{
		Operation: op,
		FileData:  fileData,
		FileId:    fileId,
	}

	*r = req

	return nil
}

func (r *Request) ParseRequestToBytes() ([]byte, error) {
	if r.FileData == nil && r.FileId == "" {
		return []byte{}, fmt.Errorf("failed to parse empty fields")
	}

	var s []byte

	if r.FileId == "" {
		s = []byte(fmt.Sprintf("%v\r\n%v\r\n\r\n", r.Operation, r.FileData))
	} else if r.FileData == nil {
		s = []byte(fmt.Sprintf("%v\r\n%v\r\n\r\n", r.Operation, r.FileId))
	}

	return s, nil
}
