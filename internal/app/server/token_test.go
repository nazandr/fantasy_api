package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToken_ParseJWT(t *testing.T) {
	conf := NewConfig()
	id := primitive.NewObjectID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, -1, 0).Unix(),
		Subject:   id.Hex(),
	})
	expToken, err := token.SignedString([]byte(conf.SignatureKey))
	if err != nil {
		assert.NoError(t, err)
	}

	testCases := []struct {
		name     string
		jwtToken string
		id       primitive.ObjectID
		err      error
	}{
		{
			name:     "expired",
			jwtToken: expToken,
			id:       id,
			err:      fmt.Errorf(""),
		},
		{
			name: "new",
			id:   primitive.NewObjectID(),
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := NewToken()
			token.AcssesToken = tc.jwtToken
			if tc.name == "new" {
				token.NewJWT(tc.id, NewConfig())
			}
			id, err := token.ParseJWT(NewConfig())
			if tc.name == "new" {
				assert.Equal(t, tc.id, id)
				assert.Nil(t, err)
				return
			}
			assert.Equal(t, tc.id, id)
			assert.NotNil(t, err)
		})
	}
}
