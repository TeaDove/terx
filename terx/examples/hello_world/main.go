package main

import (
	"os"
	"time"

	"github.com/teadove/terx/terx"
)

func echo(c *terx.Ctx) error {
	return c.Reply(c.FullText)
}

func now(c *terx.Ctx) error {
	return c.Reply(time.Now().String())
}

func main() {
	app, err := terx.New(terx.Config{Token: os.Getenv("TG_TOKEN")})
	if err != nil {
		panic(err)
	}

	app.AddHandler(echo, terx.FilterIsMessage(), terx.FilterNotCommand())
	app.AddHandler(now, terx.FilterCommand("now"))
	app.PollerRun()
}
