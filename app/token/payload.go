package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpired = errors.New("token has expired")
	ErrInvalid = errors.New("token is invalid")
)

type Payload struct {
	ID  uuid.UUID `json:"id"`
	UserName string `json:"user_name"`
	IssueAt time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (payload *Payload, err error) {
	var tokeId uuid.UUID
	tokeId , err = uuid.NewRandom()
	if err!= nil {
        return
    }

    payload = &Payload{
		ID: tokeId,
		UserName: username,
		IssueAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return
}

func (payload *Payload) Valid() (err error){
	if time.Now().After(payload.ExpiredAt) {
		err = ErrExpired
	}

	return
}