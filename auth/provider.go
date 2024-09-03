package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type Credentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthProvider interface {
	Authenticate(username, password string) (string, error)
	GetSecret() []byte
}

type JwtAuthProvider struct {
	Secret string
}

func NewAuthProvider() AuthProvider {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return NewJwtAuthProvider(secret)
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

func (p *JwtAuthProvider) GetSecret() []byte {
	return []byte(p.Secret)
}

func (p *JwtAuthProvider) generateToken(username string) (string, error) {
	now := time.Now()
	twelveHoursFromNow := now.Add(time.Hour * 12)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: now.Unix(),
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: twelveHoursFromNow.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.Secret))
}
