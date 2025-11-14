package handlers

import (
	"go-file-streamer/storage"
	. "go-file-streamer/utils"
	"net"
)

func (h *Handler) Upload(fileData *[]byte, saver storage.Saver) ([]byte, error) {

	fileId, err := saver.SaveFile(fileData)

	if err != nil {
		return []byte{}, err
	}

	res := UploadResponse{
		IsError: false,
		Message: "Received file data",
		DataId:  fileId,
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
