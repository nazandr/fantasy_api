package store_test

import (
	"os"
	"testing"
)

var database_url string

func TestMain(m *testing.M) {
	database_url = os.Getenv("DATABASE_URL")

	if database_url == "" {
		database_url = "mongodb://localhost:27017/"
	}

	os.Exit(m.Run())
}
