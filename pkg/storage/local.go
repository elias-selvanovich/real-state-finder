package storage

import (
	"io/ioutil"
)

type Storage interface {
	Load() ([]byte, error)
	Save([]byte) error
}

type storage struct {
	filepath string
}

func NewStorage(filepath string) Storage {
	return &storage{filepath: filepath}
}

func (s *storage) Load() ([]byte, error) {
	bArr, err := ioutil.ReadFile(s.filepath)
	if err != nil {
		return nil, err
	}

	return bArr, nil
}

func (s *storage) Save(data []byte) error {
	err := ioutil.WriteFile(s.filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
