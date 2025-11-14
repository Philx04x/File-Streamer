package utils

type Operation int

const (
	SAVE Operation = iota
	RETRIEVE
)

type UploadResponse struct {
	IsError bool
	Message string
	DataId  string
}

type RetrieveResponse struct {
	IsError  bool
	Message  string
	FileData *[]byte
}

type Request struct {
	Operation Operation
	FileData  []byte
	FileId    string
}
