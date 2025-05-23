# Terx
Very simple Telegram BOT framework build on top of [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api/v5), inspired by [python-telegram-bot-api](https://python-telegram-bot.org/) and [fiber](https://gofiber.io/)

## Quick start
```go
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
```

And run it with: `TG_TOKEN=you-token go run main.go`

## Features 
- Command parsers 
- Filters: IsCallback, IsMessage, Command
- Beautifully contexted zerolog logs 
- Err handling: you can send errors to bot owner, back to actor or implement ErrHandler: `func(c *Ctx, err error)`
