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

func TestNoop(t *testing.T) {
	p := Noop()
	if p != "" {
		t.Errorf("should return an empty PSK")
	}
}

func TestFromString(t *testing.T) {
	p := FromString("private")
	if p != "private" {
		t.Errorf("should return 'private' but was '%s'", p)
	}
}

func TestString(t *testing.T) {
	p := FromString("private")
	if p.String() != "private" {
		t.Errorf("should return 'private' but was '%s'", p)
	}
}
