package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct{ secret []byte }

func New(secret []byte) *JWT { return &JWT{secret: secret} }

func (j *JWT) Sign(phone string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   phone,
		"phone": phone,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JWT) Parse(tokenStr string) (string, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad signing method")
		}
		return j.secret, nil
	})
	if err != nil || !tok.Valid { return "", errors.New("invalid token") }
	claims, ok := tok.Claims.(jwt.MapClaims); if !ok { return "", errors.New("invalid claims") }
	phone, _ := claims["phone"].(string); if phone == "" { return "", errors.New("no phone") }
	return phone, nil
}
