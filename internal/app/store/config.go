package store

type Config struct {
	Database_url string `toml:"database_url"`
	DbName       string
}

func NewConfig() *Config {
	return &Config{
		Database_url: "mongodb://mongodb:27017/",
		DbName:       "fantasy",
	}
}
