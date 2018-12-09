package main

import (
	"github.com/JVecsei/fritz-guest/guestmanager"
	"github.com/JVecsei/fritz-guest/session"
	"github.com/kr/pretty"
)

func main() {
	s, err := session.NewSessionByPassword(url, password)
	pretty.Println(s)
	pretty.Println(err)

	g, err := guestmanager.NewGuestManager(s)
	pretty.Println(err)
	err = g.TurnOn()
	pretty.Println(err)

}
