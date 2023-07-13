package repository

import (
	"errors"
	"on-air/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type payload struct {
	UserID    int       `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}

func newPayload(userID int, duration time.Duration) *payload {
	payload := &payload{
		UserID:    userID,
		ExpiredAt: time.Now().Add(duration),
	}

	return payload
}

func (payload *payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token expired")
	}

	return nil
}

func CreateToken(cfg *config.JWT, userID int) (string, error) {
	payload := newPayload(userID, cfg.ExpiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(cfg.SecretKey))
}

var ErrInvalidToken = errors.New("invalid token")

func VerifyToken(cfg *config.JWT, token string) (*payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("")
		}
		return []byte(cfg.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
