package infrastructure

import (
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/db"
	"schbot/internal/handles"
	"schbot/internal/interfaces"
)

var env *environment

type environment struct {
	logger   *zap.SugaredLogger
	cfg      *config.Config
	client   *http.Client
	bot      *tb.Bot
	dbClient interfaces.IDBHandler
	dbMagic  interfaces.IDbMagic
}

type IInjector interface {
	InjectHandles() handles.Handles
}

func (e *environment) InjectHandles() handles.Handles {
	return handles.Handles{
		Log:       e.logger,
		Bot:       e.bot,
		Config:    e.cfg,
		Client:    e.client,
		DBHandler: e.dbClient,
		DbMagic: &db.DMagic{
			Log:    e.logger,
			Config: e.cfg,
		},
	}
}

func Injector(log *zap.SugaredLogger, bot *tb.Bot, cfg *config.Config) (IInjector, error) {
	client, err := InitPostgresClient(cfg)
	if err != nil {
		log.Fatal("injector: db client init error")
		return nil, err
	}

	env = &environment{
		logger:   log,
		cfg:      cfg,
		client:   http.DefaultClient,
		bot:      bot,
		dbClient: client,
	}

	return env, nil
}
