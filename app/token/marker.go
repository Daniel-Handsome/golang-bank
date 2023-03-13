package token

import "time"

type Marker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload ,error)
	// verify檢查軟件是否滿足規範的過程
	VerifyToken(token  string) (*Payload, error)
}