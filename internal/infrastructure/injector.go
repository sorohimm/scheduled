package infrastructure

import (
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/handles"
)

var env *environment

type environment struct {
	logger *zap.SugaredLogger
	cfg    *config.Config
	client *http.Client
	bot    *tb.Bot
}

type IInjector interface {
	InjectHandles() handles.Handles
}

func (e *environment) InjectHandles() handles.Handles {
	return handles.Handles{
		Log:    e.logger,
		Bot:    e.bot,
		Config: e.cfg,
		Client: e.client,
	}
}

func Injector(log *zap.SugaredLogger, bot *tb.Bot, cfg *config.Config) (IInjector, error) {
	env = &environment{
		logger: log,
		cfg:    cfg,
		client: http.DefaultClient,
		bot:    bot,
	}

	return env, nil
}
