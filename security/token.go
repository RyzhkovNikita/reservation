package security

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type TokenType uint

var ExpiredTokenError = errors.New("token is expired")
var InvalidTokenError = errors.New("token is expired")

const (
	Access TokenType = iota
	Refresh
)

type TokenManager interface {
	CreateToken(userId uint64, duration time.Duration, tokenType TokenType) (string, error)
	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}

type JWTTokenManager struct {
	SecretKey string
}

type Payload struct {
	UserId    uint64
	ExpiresAt time.Time
	TokenType TokenType
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ExpiredTokenError
	}
	return nil
}

func (J JWTTokenManager) CreateToken(
	userId uint64,
	duration time.Duration,
	tokenType TokenType,
) (string, error) {
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Payload{
		UserId:    userId,
		ExpiresAt: time.Now().Add(duration),
		TokenType: tokenType,
	})
	return jwttoken.SignedString([]byte(J.SecretKey))
}

func (J JWTTokenManager) VerifyToken(
	token string,
	tokenType TokenType,
) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenError
		}
		return []byte(J.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredTokenError) {
			return nil, ExpiredTokenError
		}
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenError
	}
	if payload.TokenType != tokenType {
		return nil, InvalidTokenError
	}
	return payload, nil
}
