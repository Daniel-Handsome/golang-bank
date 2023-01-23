package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

type JwtMark  struct {
	secretKey string
}

func NewJwMarkt(secretKey string)  (marker Marker, err error) {
	if len(secretKey) < minSecretKeySize {
		err  = fmt.Errorf("invalid secret key must  be at  least  %d char", minSecretKeySize)
		return 
	}

	marker = &JwtMark{secretKey}

	return
}

func (jwtMark *JwtMark) CreateToken(username string, duration time.Duration) (string, error) {
	paload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	//對稱加密
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, paload)
	return claim.SignedString([]byte(jwtMark.secretKey))
}

func (jwtMark *JwtMark) VerifyToken(token  string) (*Payload, error) {
	keyFunc :=  func(token *jwt.Token) (interface{}, error) {
		// valida method 
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalid
		}

		return []byte(jwtMark.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if validatetionerr, ok := err.(*jwt.ValidationError); ok {
			if errors.Is(validatetionerr.Inner, ErrExpired) {
				return nil, ErrExpired
			}
		}
		return nil, ErrInvalid
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalid
	}

	return payload, nil
}