package utils

import (
	"math/rand"
	"time"
)

const alphber string = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixMicro())
}

func RandInt(min, max int) int {
	return rand.Intn(max - min) + min
}

func RandString(n int)  string{
	// 1
	// var sb strings.Builder
	// for i := 0; i < n; i++ {
	// 	sb.WriteByte(alphber[rand.Intn(len(alphber))])
	// }

	// return sb.String()

	// 2
	str := make([]rune, n)

	for i := range str {
		str[i] = rune(alphber[rand.Intn(len(alphber))])
	}

	return string(str)
}

func RandCurrency() string {
	currencies := []string{"USD", "EUR", "GBP"}
	return currencies[rand.Intn(len(currencies))]
}
