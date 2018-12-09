package psk

import (
	"math/rand"
	"time"
)

//Psk used to configure the fritzbox
type Psk string

//FromString returns psk from string
func FromString(s string) Psk {
	return Psk(s)
}

//Random returns random psk with specific length
func Random(length int) Psk {
	return Psk(generateRandomString(length))
}

//Noop returns an empty PSK
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
