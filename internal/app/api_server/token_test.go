package api_server

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, -1, 0).Unix(),
		Subject:   "1",
	})
	expToken, err := token.SignedString([]byte(conf.SignatureKey))
	if err != nil {
		assert.NoError(t, err)
	}

	testCases := []struct {
		name     string
		jwtToken string
		err      error
	}{
		{
			name:     "expired",
			jwtToken: expToken,
			err:      fmt.Errorf(""),
		},
		{
			name: "new",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newId := primitive.NewObjectID()
			token := NewToken()
			token.AcssesToken = tc.jwtToken
			if tc.name == "new" {
				token.NewJWT(newId, NewConfig())
			}
			id, err := token.ParseJWT(NewConfig())
			if tc.name == "new" {
				assert.Equal(t, newId, id)
				assert.Nil(t, err)
				return
			}
			assert.NotNil(t, err)
		})
	}
}
