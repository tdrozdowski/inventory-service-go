package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtAuthProvider_Authenticate(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			username: "foo",
			password: "bar",
			wantErr:  false,
		},
		{
			name:     "Invalid password",
			username: "foo",
			password: "baz",
			wantErr:  true,
		},
		{
			name:     "No username",
			username: "",
			password: "bar",
			wantErr:  true,
		},
		{
			name:     "No password",
			username: "foo",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewJwtAuthProvider("dummy_secret")
			token, err := provider.Authenticate(tt.username, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				decodedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte("dummy_secret"), nil
				})
				assert.NotNil(t, decodedToken)
				assert.True(t, decodedToken.Valid)
				assert.Equal(t, "foo", decodedToken.Claims.(jwt.MapClaims)["username"])
				assert.NoError(t, err)
			}
		})
	}
}
