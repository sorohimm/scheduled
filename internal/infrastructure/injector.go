package infrastructure

import (
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/handles_controllers"
	"schbot/internal/handles_repos"
	"schbot/internal/handles_services"
	"schbot/internal/interfaces"
	"schbot/internal/schedule_maker"
)

var env *environment

type environment struct {
	logger   *zap.SugaredLogger
	cfg      *config.Config
	client   *http.Client
	bot      *tb.Bot
	dbClient interfaces.IDBHandler
}

type IInjector interface {
	InjectHandles() handles_controllers.HandlesController
}

func (e *environment) InjectHandles() handles_controllers.HandlesController {
	return handles_controllers.HandlesController{
		Log:    e.logger,
		Bot:    e.bot,
		Config: e.cfg,
		HandlesService: &handles_services.HandleService{
			Log:           e.logger,
			Config:        e.cfg,
			Client:        e.client,
			DBHandler:     e.dbClient,
			ScheduleMaker: &schedule_maker.ScheduleMaker{},
			HandlesRepo: &handles_repos.HandlesRepo{
				Log:    e.logger,
				Config: e.cfg,
				Client: e.client,
			},
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
