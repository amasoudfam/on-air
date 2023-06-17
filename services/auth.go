package services

import (
	"errors"
	"on-air/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	UserID    int       `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}

type TokenMaker interface {
	CreateToken(cfg *config.Config, userID int) (string, error)

	VerifyToken(token string) (*Payload, error)
}

func NewPayload(userID int, duration time.Duration) *Payload {
	payload := &Payload{
		UserID:    userID,
		ExpiredAt: time.Now().Add(duration),
	}
	return payload
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("payload expired")
	}
	return nil
}

func CreateToken(cfg *config.Config, userID int) (string, error) {
	payload := NewPayload(userID, time.Duration(cfg.Auth.LifeTime*int(time.Minute)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}

var ErrInvalidToken = errors.New("invalid token")

func VerifyToken(cfg *config.Config, token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("")
		}
		return []byte(cfg.Auth.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
