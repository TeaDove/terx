# Terx
Very simple Telegram BOT framework build on top of github.com/go-telegram-bot-api/telegram-bot-api/v5

## Quick start
```go

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
```

And run it with: `TG_TOKEN=you-token go run main.go`

## Features 
- Command parsers 
- Filters: IsCallback, IsMessage, Command
- Beautifully contexted zerolog logs 
