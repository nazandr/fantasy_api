package store

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestStore(t *testing.T, database_url string) (*Store, func(*mongo.Collection)) {
	t.Helper()

	config := NewConfig()
	config.Database_url = database_url

	s := New(config)
	if err := s.Connect(); err != nil {
		t.Fatal()
	}

	return s, func(c *mongo.Collection) {
		if err := c.Drop(s.context); err != nil {
			t.Fatal(err)
		}
		s.Close()
	}

}
