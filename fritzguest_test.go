package fritzguest

import (
	"github.com/JVecsei/fritz-guest/psk"
	"log"
	"github.com/JVecsei/fritz-guest/guestmanager"
	"github.com/JVecsei/fritz-guest/session"
)

//Turn on guest access with currently set PSK
func Example_turnOn() {
	s, err := session.NewSessionByUsernamePassword("http://fritz.box", "uName", "passwd")
	if err != nil {
		handleError(err)
	}
	g, err := guestmanager.NewGuestManager(s)
	if err != nil {
		handleError(err)
	}
	err = g.TurnOn()
	if err != nil {
		handleError(err)
	}
}

//Turn on guest access with a new given PSK
func Example_turnOnWithSpecificPsk() {
	s, err := session.NewSessionByUsernamePassword("http://fritz.box", "uName", "passwd")
	if err != nil {
		handleError(err)
	}
	g, err := guestmanager.NewGuestManager(s)
	if err != nil {
		handleError(err)
	}
	err = g.TurnOnWithPsk(psk.FromString("newPsk"))
	if err != nil {
		handleError(err)
	}
}

//Turn on guest access with a random PSK
func Example_turnOnWithRandomPsk() {
	s, err := session.NewSessionByUsernamePassword("http://fritz.box", "uName", "passwd")
	if err != nil {
		handleError(err)
	}
	g, err := guestmanager.NewGuestManager(s)
	if err != nil {
		handleError(err)
	}
	err = g.TurnOnWithPsk(psk.Random(10))
	if err != nil {
		handleError(err)
	}
}

//Turn off guest access
func Example_turnOff() {
	s, err := session.NewSessionByUsernamePassword("http://fritz.box", "uName", "passwd")
	if err != nil {
		handleError(err)
	}
	g, err := guestmanager.NewGuestManager(s)
	if err != nil {
		handleError(err)
	}
	err = g.TurnOff()
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	log.Fatalf("error: %v", err)
}