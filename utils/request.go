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
