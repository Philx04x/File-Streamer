package utils

import (
	"bytes"
	"fmt"
	"strconv"
)

func (r *RetrieveResponse) ParseResponseToStruct(resBytes []byte) error {

	// Response:
	// true\r\n
	// blablabla\r\n
	// data\r\n
	// \r\n
	lines := bytes.Split(resBytes, []byte("\r\n"))

	line := string(lines[0])
	msg := string(lines[1])
	fileData := lines[2]

	isError, _ := strconv.ParseBool(line)

	res := RetrieveResponse{
		IsError:  isError,
		Message:  msg,
		FileData: &fileData,
	}

	*r = res

	return nil
}

func (r *RetrieveResponse) ParseResponseToBytes() ([]byte, error) {
	if r.Message == "" {
		return []byte{}, fmt.Errorf("failed to parse empty field")
	}

	s := []byte(fmt.Sprintf("%v\r\n%s\r\n%s\r\n\r\n", r.IsError, r.Message, *(r.FileData)))

	return s, nil
}
