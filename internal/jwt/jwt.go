package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidated = errors.New("jwt has been invalidated")
)

type Jwt struct {
	secret string
}

func New(secret string) *Jwt {
	return &Jwt{
		secret: secret,
	}
}

type claims struct {
	jwt.StandardClaims
}

func (j *Jwt) New(expires time.Duration) (string, error) {
	exp := time.Now().Add(expires)

	claims := claims{
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *Jwt) Verify(ctx context.Context, token string) error {
	_, err := j.getClaims(token)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jwt) getClaims(token string) (*claims, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
