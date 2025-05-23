package terx

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
)

type Ctx struct {
	Context context.Context
	Terx    *Terx

	// Text
	// Update's text without Command
	Text string

	// FullText
	// Update.Message.Text
	FullText string

	Command string

	Update tgbotapi.Update

	// SentFrom
	// Can be nil!
	SentFrom *tgbotapi.User

	// Chat
	// Can be nil!
	Chat *tgbotapi.Chat
}

func (r *Ctx) addLogCtx() {
	if r.Chat != nil && r.Chat.Title != "" {
		r.Context = logger_utils.WithValue(r.Context, "in", r.Chat.Title)
	}

	if r.Text != "" {
		r.Context = logger_utils.WithValue(r.Context, "text", redact_utils.Trim(r.Text))
	}

	if r.SentFrom != nil {
		r.Context = logger_utils.WithValue(r.Context, "from", r.SentFrom.String())
	}

	if r.Command != "" {
		r.Context = logger_utils.WithValue(r.Context, "command", r.Command)
	}
}

func (r *Terx) makeCtx(ctx context.Context, update *tgbotapi.Update) *Ctx {
	c := Ctx{
		Terx:     r,
		Update:   *update,
		Chat:     update.FromChat(),
		SentFrom: update.SentFrom(),
	}

	if update.Message != nil {
		c.FullText = update.Message.Text
	}

	c.Context = ctx

	inChat := c.SentFrom != nil && c.Chat != nil && c.SentFrom.ID != c.Chat.ID
	c.Command, c.Text = extractCommandAndText(c.FullText, r.Bot.Self.UserName, inChat)
	c.Text = strings.TrimSpace(c.Text)

	c.addLogCtx()

	return &c
}

func (r *Ctx) Log() *zerolog.Logger {
	return zerolog.Ctx(r.Context)
}

func (r *Ctx) LogWithUpdate() *zerolog.Logger {
	logger := zerolog.Ctx(r.Context).With().Interface("update", r.Update).Logger()
	return &logger
}
