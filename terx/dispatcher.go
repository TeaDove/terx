package terx

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/must_utils"
	"github.com/teadove/teasutils/utils/reflect_utils"
)

type (
	ProcessorFunc func(r *Ctx) error
	FilterFunc    func(r *Ctx) bool
)

type Handler struct {
	Filters   []FilterFunc
	Processor ProcessorFunc
}

func (r *Terx) PollerRun(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	// TODO move to settings
	u.Timeout = 10
	updates := r.Bot.GetUpdatesChan(u)

	zerolog.Ctx(ctx).
		Info().
		Interface("handlers", len(r.Handlers)).
		Msg("bot.polling.started")

	var wg sync.WaitGroup

	for update := range updates {
		wg.Add(1)
		// TODO реализовать разные методы параллелизмма

		go r.processUpdate(ctx, &wg, &update)
	}

	wg.Wait()
}

func (r *Terx) processUpdate(ctx context.Context, wg *sync.WaitGroup, update *tgbotapi.Update) {
	defer func() {
		err := must_utils.AnyToErr(recover())
		if err == nil {
			return
		}

		zerolog.Ctx(ctx).
			Error().
			Stack().Err(err).
			Interface("update", update).
			Msg("panic.in.process.update")
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer wg.Done()

	c := r.makeCtx(ctx, update)

	for _, handler := range r.Handlers {
		ok := true

		for _, filter := range handler.Filters {
			if !filter(c) {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		err := handler.Processor(c)
		if err != nil {
			r.handleError(c, err, handler)
			continue
		}

		c.Log().Debug().
			Str("processor", reflect_utils.GetFunctionName(handler.Processor)).
			Msg("handler.processed")
	}

	c.Log().Debug().
		Msg("update.processed")
}

func (r *Terx) handleError(c *Ctx, err error, handler Handler) {
	err = errors.Wrap(err, "failed to process handler")

	c.LogWithUpdate().Error().
		Stack().Err(err).
		Str("processor", reflect_utils.GetFunctionName(handler.Processor)).
		Msg("failed.to.process.handler")

	if r.replyWithErr {
		c.TryReplyOnErr(err)
	}

	if r.sendErrToOwner {
		msgReq := tgbotapi.NewMessage(
			r.ownerUserID,
			fmt.Sprintf(`Error occurred\n\nerr: %s\n\nupdate: %+v`, err.Error(), c.Update),
		)

		_, err = r.Bot.Send(msgReq)
		if err != nil {
			c.LogWithUpdate().
				Error().
				Stack().Err(err).
				Msg("failed.to.send.message.on.error")
		}
	}

	if r.errHandler != nil {
		r.errHandler(c, err)
	}
}
