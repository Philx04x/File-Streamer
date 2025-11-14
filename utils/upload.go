package utils

import (
	"bytes"
	"fmt"
	"strconv"
)

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
