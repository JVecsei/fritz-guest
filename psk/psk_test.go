package psk

import (
	"testing"
)

func TestRandom(t *testing.T) {
	length := 20
	p := Random(length)
	if len(p) != 20 {
		t.Errorf("invalid length of random psk %d != %d", len(p), length)
	}
}
