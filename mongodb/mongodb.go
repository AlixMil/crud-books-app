package mongodb

type MongoDB struct{}

func New(username, password, dsn string) (*MongoDB, error) {
	return &MongoDB{}, nil
}
