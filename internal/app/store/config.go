package store

type Config struct {
	Database_url string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
