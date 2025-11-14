package handlers

import (
	. "go-file-streamer/utils"
	"net"
)

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
		DataId:  "",
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
