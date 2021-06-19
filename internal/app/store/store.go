package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	config                *Config
	context               context.Context
	db                    *mongo.Database
	userCollection        *UserCollection
	playerCardsCollection *PlayerCardsCollection
	buffsCollection       *BuffsCollection
	matchCollection       *MatchCollection
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

	s.db = client.Database(s.config.DbName)
	return nil
}

func (s *Store) Close() {
	s.db.Client().Disconnect(s.context)
}
func (s *Store) DropDb() {
	s.db.Drop(s.context)
}

func (s *Store) User() *UserCollection {
	if s.userCollection != nil {
		return s.userCollection
	}

	s.userCollection = &UserCollection{
		Store:      s,
		Collection: s.db.Collection("users"),
	}
	return s.userCollection
}

func (s *Store) PlayerCards() *PlayerCardsCollection {
	if s.playerCardsCollection != nil {
		return s.playerCardsCollection
	}

	s.playerCardsCollection = &PlayerCardsCollection{
		Store:      s,
		Collection: s.db.Collection("player_cards"),
	}
	return s.playerCardsCollection
}

func (s *Store) Buffs() *BuffsCollection {
	if s.buffsCollection != nil {
		return s.buffsCollection
	}

	s.buffsCollection = &BuffsCollection{
		Store:      s,
		Collection: s.db.Collection("card_buffs"),
	}
	return s.buffsCollection
}

func (s *Store) Matches() *MatchCollection {
	if s.matchCollection != nil {
		return s.matchCollection
	}

	s.matchCollection = &MatchCollection{
		Store:      s,
		Collection: s.db.Collection("matches"),
	}
	return s.matchCollection
}
