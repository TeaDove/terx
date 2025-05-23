package terx

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Terx struct {
	Bot      *tgbotapi.BotAPI
	Handlers []Handler
	LogLevel zerolog.Level

	errHandler     func(c *Ctx, err error)
	replyWithErr   bool
	sendErrToOwner bool
	ownerUserID    int64
}

type Config struct {
	Token    string
	LogLevel zerolog.Level

	ErrHandler     func(c *Ctx, err error)
	ReplyWithErr   bool
	SendErrToOwner bool
	OwnerUserID    int64
}

func New(config *Config) (*Terx, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bot client")
	}

	terx := &Terx{
		Bot:            bot,
		LogLevel:       config.LogLevel,
		Handlers:       make([]Handler, 0),
		errHandler:     config.ErrHandler,
		replyWithErr:   config.ReplyWithErr,
		sendErrToOwner: config.SendErrToOwner,
		ownerUserID:    config.OwnerUserID,
	}

	if config.SendErrToOwner && config.OwnerUserID == 0 {
		return nil, errors.New("OwnerUserID must be set if SendErrToOwner is true")
	}

	return terx, nil
}

func (r *Terx) AddHandler(processor ProcessorFunc, filters ...FilterFunc) {
	if processor == nil {
		panic("processor cannot be nil")
	}

	r.Handlers = append(r.Handlers, Handler{
		Filters:   filters,
		Processor: processor,
	})
}
