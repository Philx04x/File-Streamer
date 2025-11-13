package handlers

import "errors"

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

var ErrFailedParseToBytes error = errors.New("failed to parse data to bytes")
