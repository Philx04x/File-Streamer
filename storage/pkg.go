package storage

type Saver struct {
	Service ISaver
}

type ISaver interface {
	SaveFile(*[]byte) (*string, error)
	RetrieveFile(string) (*[]byte, error)
	BuildUpCache() error
}

func NewSaverService(path string) Saver {
	ser := NewFileSaver(path)

	return Saver{
		Service: ser,
	}
}

func (s *Saver) SaveFile(fileData *[]byte) (*string, error) {
	return s.Service.SaveFile(fileData)
}

func (s *Saver) RetrieveFile(fileId string) (*[]byte, error) {
	return s.Service.RetrieveFile(fileId)
}

func (s *Saver) BuildUpCache() error {
	return s.Service.BuildUpCache()
}
