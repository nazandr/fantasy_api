package store

type Config struct {
	dbURL string
}

func NewConfig() *Config {
	return &Config{}
}

type Store struct {
	config *Config
}

func (c *Config) New() *Store {
	return &Store{
		config: c,
	}
}

func (s *Store) Connect() error {
	return nil
}

func (s *Store) Close() {
}
