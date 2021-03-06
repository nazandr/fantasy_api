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

	u := models.TestUser(t)
	err := s.User().Create(u)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUser_collection_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, database_url)
	s.User().Collection = s.User().Collection.Database().Collection("user_test")

	defer teardown(s.User().Collection)

	u := models.TestUser(t)
	_, err := s.User().FindByEmail(u.Email)
	assert.Error(t, err)

	s.User().Create(u)

	u, err = s.User().FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUser_collection_Find(t *testing.T) {
	s, teardown := store.TestStore(t, database_url)
	s.User().Collection = s.User().Collection.Database().Collection("user_test")

	defer teardown(s.User().Collection)

	u := models.TestUser(t)

	s.User().Create(u)

	fu, err := s.User().Find(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fu)
}

func TestUser_collerction_ReplaseUser(t *testing.T) {
	s, teardown := store.TestStore(t, database_url)
	s.User().Collection = s.User().Collection.Database().Collection("user_test")

	defer teardown(s.User().Collection)

	u := models.TestUser(t)

	s.User().Create(u)

	uu := u
	uu.Packs.Common = 5

	err := s.User().ReplaseUser(uu)
	assert.NoError(t, err)

	fu, err := s.User().Find(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, uu.Packs.Common, fu.Packs.Common)
}

func TestUser_collection_UpdateRefreshToekn(t *testing.T) {
	s, teardown := store.TestStore(t, database_url)
	s.User().Collection = s.User().Collection.Database().Collection("user_test")
	defer teardown(s.User().Collection)

	u := models.TestUser(t)
	s.User().Create(u)
	err := s.User().UpdateRefreshToken(u.ID, "asdfsadf", 30)
	assert.NoError(t, err)
}
