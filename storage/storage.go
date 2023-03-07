package storage

type StorageActions interface {
	GetServer() string
	UploadFileSummary() map[string]string
}

type Storage struct {
	actions StorageActions
}

func New(service StorageActions) (*Storage, error) {
	return &Storage{actions: service}, nil
}
