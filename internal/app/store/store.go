package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	dbURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}

type Store struct {
	config  *Config
	context context.Context
	db      *mongo.Client
}

func New(c *Config) *Store {
	return &Store{
		config:  c,
		context: context.TODO(),
	}
}

func (s *Store) Connect() error {
	clientOptions := options.Client().ApplyURI(s.config.dbURL)
	client, err := mongo.Connect(s.context, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(s.context, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() {
	s.db.Disconnect(s.context)
}
