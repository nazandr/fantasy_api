package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	clientOptions := options.Client().ApplyURI(s.config.Database_url)
	client, err := mongo.Connect(s.context, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(s.context, nil)
	if err != nil {
		return err
	}

	s.db = client
	return nil
}

func (s *Store) Close() {
	s.db.Disconnect(s.context)
}
