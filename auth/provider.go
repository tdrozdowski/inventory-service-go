package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthProvider interface {
	Authenticate(username, password string) (string, error)
}

type JwtAuthProvider struct {
	Secret string
}

func NewJwtAuthProvider(secret string) *JwtAuthProvider {
	return &JwtAuthProvider{
		Secret: secret,
	}
}

func (p *JwtAuthProvider) Authenticate(username, password string) (string, error) {
	// TODO - lookup from db eventually
	if username == "foo" && password == "bar" {
		return p.generateToken(username)
	} else {
		return "", errors.New("Invalid credentials")
	}
}

func (p *JwtAuthProvider) generateToken(username string) (string, error) {
	standardExpirationTime := time.Now().Add(time.Hour * 12).Unix()
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: standardExpirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.Secret))
}
