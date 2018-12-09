package psk

import (
	"math/rand"
	"time"
)

type Psk string

func FromString(s string) Psk {
	return Psk(s)
}

func Random(length int) Psk {
	return Psk(generateRandomString(length))
}

func Noop() Psk {
	return Psk("")
}

func (p Psk) String() string {
	return string(p)
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!?_123456789"
	psk := make([]byte, length)
	for i := range psk {
		psk[i] = chars[rand.Intn(len(chars))]
	}
	return string(psk)
}
