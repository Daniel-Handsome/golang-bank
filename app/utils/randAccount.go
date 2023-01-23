package utils

import (
	"fmt"
	"math/rand"
)

func RandOwner() string {
	return RandString(5)
}

func RandBalance() int64 {
	return int64(RandInt(10, 1000))
}

func RandCurrency() string {
	return currencies[rand.Intn(len(currencies))]
}

func RandEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandString(6))
}
