package store

type Config struct {
	dbURL string
}

type Store struct {
	config *Config
}

func (c *Config) New() *Store {
	return &Store{
		config: c,
	}
}
