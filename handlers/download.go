package handlers

import (
	"go-file-streamer/storage"
	. "go-file-streamer/utils"
)

func (h *Handler) Download(fileId string, saver storage.Saver) ([]byte, error) {
	fileData, err := saver.RetrieveFile(fileId)

	if err != nil {
		return []byte{}, err
	}

	res := RetrieveResponse{
		IsError:  false,
		Message:  "Retrieved file data",
		FileData: fileData,
	}

	b, err := res.ParseResponseToBytes()

	if err != nil {
		return []byte{}, ErrFailedParseToBytes
	}

	return b, nil
}
