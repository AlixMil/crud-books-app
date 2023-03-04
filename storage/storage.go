package storage

import "fmt"

type TokenGetter interface {
	Token() string
}

type StorageActions interface {
	SaveFile(file []byte, title string) (tgtr TokenGetter, err error)
}

type Storage struct {
	actions StorageActions
}

func (s *Storage) SaveFile(file []byte, title string) (string, error) {
	tgtr, err := s.actions.SaveFile(file, title)
	if err != nil {
		return "", fmt.Errorf("save file error: %v", err)
	}
	return tgtr.Token(), nil
}

func New(fileSaver StorageActions) (*Storage, error) {

	return &Storage{actions: fileSaver}, nil
}
