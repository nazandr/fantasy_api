package models_test

import (
	"testing"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := models.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncriptedPassword)
}
