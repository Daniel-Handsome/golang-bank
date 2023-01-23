package token

import (
	"testing"
	"time"

	"github.com/daniel/master-golang/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTMarker(t *testing.T) {
	jwtMark, err := NewJwMarkt(utils.RandString(32))
	require.NoError(t, err)

	username := utils.RandOwner()
	duration  := time.Minute

	issuedAt := time.Now()
	expired := issuedAt.Add(duration)

	token, err := jwtMark.CreateToken(username, duration)
	require.NoError(t, err)
    require.NotEmpty(t, token)

	paload, err := jwtMark.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, paload)

	// uuid
	require.NotZero(t, paload.ID)
	require.Equal(t, paload.UserName, username)
	require.WithinDuration(t, paload.IssueAt, issuedAt, time.Second)
	require.WithinDuration(t, paload.ExpiredAt, expired, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	jwtMark, err := NewJwMarkt(utils.RandString(32))
    require.NoError(t, err)

	username := utils.RandOwner()
    duration  := -time.Minute

	token, err := jwtMark.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMark.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpired.Error())

	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(utils.RandOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	//去看newClaim註冊的metod去找他的sign方法 就知道他是哪個type 
	// 這只能用來test 
	toke, err :=jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	jwtMark, err := NewJwMarkt(utils.RandString(32))
	payload, err = jwtMark.VerifyToken(toke)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalid.Error())
	require.Nil(t, payload)
}