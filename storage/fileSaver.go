package storage

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
)

func NewFileSaver(path string) ISaver {
	return &FileSaver{
		filesCache: make(map[string]*[]byte, 10),
		path:       path,
	}
}

type FileSaver struct {
	fileMu     sync.RWMutex
	filesCache map[string]*[]byte
	path       string
}

func (s *FileSaver) SaveFile(fileData *[]byte) (string, error) {
	fileId := uuid.NewString()

	err := os.WriteFile(fmt.Sprintf("%s/%s.bin", s.path, fileId), *fileData, 0644)
	if err != nil {
		return "", ErrFileWrite
	}

	s.fileMu.Lock()
	s.filesCache[fileId] = fileData
	s.fileMu.Unlock()

	return fileId, nil
}

func (s *FileSaver) RetrieveFile(fileId string) (*[]byte, error) {
	cacheData := s.filesCache[fileId]

	if cacheData != nil {
		return cacheData, nil
	}

	fileReader, err := os.Open(fmt.Sprintf("%s/%s.bin", s.path, fileId))
	if err != nil {
		return nil, ErrFileOpen
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)

	if err != nil {
		return nil, ErrFileReader
	}

	return &fileData, nil
}

func (s *FileSaver) BuildUpCache() error {
	dir, err := os.ReadDir(s.path)

	if err != nil {
		return ErrFileDirReader
	}

	for _, file := range dir {

		fName := file.Name()
		fName = fName[:len(fName)-4]

		if s.filesCache[fName] != nil {
			continue
		}

		fileReader, err := os.Open(fmt.Sprintf("%s/%s.bin", s.path, fName))

		if err != nil {
			return ErrFileOpen
		}

		defer fileReader.Close()

		fileData, err := io.ReadAll(fileReader)

		if err != nil {
			return ErrFileReader
		}

		s.filesCache[fName] = &fileData

	}

	return nil
}
