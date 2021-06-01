package api_server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToken_ParseJWT(t *testing.T) {
	testCases := []struct {
		name     string
		jwtToken string
		err      error
	}{
		{
			name:     "expired",
			jwtToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjI0OTU0NDMsInN1YiI6Ik9iamVjdElEKFwiNjBiM2U0Y2QyMTYzNjhiNjY4OGI3NWU3XCIpIn0.nZ_jDsAAIxnQJnQlkvzOU-CvA6w_cFN31tHG0uV51Qo",
			err:      fmt.Errorf(""),
		},
		{
			name: "new",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := NewToken()
			token.AcssesToken = tc.jwtToken
			if tc.name == "new" {
				token.NewJWT(primitive.NewObjectID(), NewConfig())
			}
			_, err := token.ParseJWT(NewConfig())
			if tc.name == "new" {
				assert.Nil(t, err)
				return
			}
			assert.NotNil(t, err)
		})
	}
}
