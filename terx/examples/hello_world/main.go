package main

import (
	"github.com/teadove/terx/terx"
	"os"
)

func main() {
	app, err := terx.New(terx.Config{Token: os.Getenv("TG_TOKEN")})
	if err != nil {
		panic(err)
	}

	app.AddHandler(func(r *terx.Ctx) error { return r.Reply("hi!") }, terx.FilterIsMessage())
	app.PollerRun()
}
