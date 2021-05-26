package store_test

import (
	"testing"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/nazandr/fantasy_api/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestUser_collection_Create(t *testing.T) {
	s, teardown := store.TestStore(t, database_url)
	s.User().Collection = s.User().Collection.Database().Collection("user_test")

	defer teardown(s.User().Collection)

	u, err := s.User().Create(&models.User{
		Email: "user@example.com",
	})
	t.Log(s.User())

	assert.NoError(t, err)
	assert.NotNil(t, u)
}
