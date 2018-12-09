package main

import (
	"os"

	"github.com/JVecsei/fritz-guest/guestmanager"
	"github.com/JVecsei/fritz-guest/session"
	"github.com/kr/pretty"
)

func main() {
	url := os.Getenv("FBURL")
	password := os.Getenv("FBPASSWORD")
	s, err := session.NewSessionByPassword(url, password)
	pretty.Println(s)
	pretty.Println(err)

	g, err := guestmanager.NewGuestManager(s)
	pretty.Println(err)
	err = g.TurnOn()
	pretty.Println(err)

}
