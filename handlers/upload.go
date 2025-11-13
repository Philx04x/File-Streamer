package handlers

import (
	. "go-file-streamer/utils"
	"net"
)

func (h *Handler) Upload() ([]byte, error) {

	uuid := "9226334b-83d6-4b47-a8e9-4b5efcb59e3c"

	res := UploadResponse{
		IsError: false,
		Message: "Received file data",
		DataId:  &uuid,
	}

	b, err := res.ParseResponseToBytes()

	if err != nil {
		return []byte{}, ErrFailedParseToBytes
	}

	return b, nil
}

func (h *Handler) ResponseWriter(conn net.Conn, b []byte) {

	conn.Write(b)

	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}

	conn.Close()
}

func (h *Handler) ErrorResponseWriter(conn net.Conn, errorMsg string) {
	r := UploadResponse{
		IsError: true,
		Message: errorMsg,
		DataId:  nil,
	}

	b, err := r.ParseResponseToBytes()

	if err != nil {
		panic("failed to parse error response")
	}

	conn.Write(b)

	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}

	conn.Close()

}
