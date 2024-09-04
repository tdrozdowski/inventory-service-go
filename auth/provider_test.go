package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestNewAuthProvider(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   bool
	}{
		{
			name:   "Valid JWT Secret",
			secret: "dummy_secret",
			want:   true,
		},
		{
			name:   "Empty JWT Secret",
			secret: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("JWT_SECRET", tt.secret)
			defer func() {
				_ = os.Unsetenv("JWT_SECRET")
			}()

			var authProvider AuthProvider
			if tt.want {
				authProvider = NewAuthProvider()
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Log("Recovered from panic as expected")
					} else {
						t.Error("Expected a panic for empty JWT_SECRET but did not occur")
					}
				}()
				authProvider = NewAuthProvider()
			}

			if authProvider != nil {
				assert.NotNil(t, authProvider.GetSecret())
			}
		})
	}
}

func TestJwtAuthProvider_GetSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   []byte
	}{
		{
			name:   "Valid JWT Secret",
			secret: "dummy_secret",
			want:   []byte("dummy_secret"),
		},
		{
			name:   "Empty JWT Secret",
			secret: "",
			want:   []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewJwtAuthProvider(tt.secret)
			got := provider.GetSecret()
			assert.Equal(t, tt.want, got)
		})
	}
}
