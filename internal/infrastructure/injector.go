package infrastructure

import (
	"context"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/db"
	"schbot/internal/handles"
	"schbot/internal/interfaces"

	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
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

func CreateTable(log *zap.SugaredLogger, client interface{interfaces.IDBHandler}) error {
  conn, err := client.AcquireConn(context.Background())
  if err != nil {
    log.Fatal("unable to create table in db")
    return err
  }
  defer conn.Release()

  const CreateExtensionStatement = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
  const CreateTableStatement = `CREATE TABLE IF NOT EXISTS chats (
    chat_id bigint,
    group_name text,
    group_uuid UUID,
  );`

  err = conn.QueryRow(context.Background(), CreateExtensionStatement).Scan()
  if err != nil {
    log.Fatal("unable to create table in db 2")
    return err
  }
  err = conn.QueryRow(context.Background(), CreateTableStatement).Scan()
  if err != nil {
    log.Fatal("unable to create table in db 3")
    return err
  }
  return nil
}

func Injector(log *zap.SugaredLogger, bot *tb.Bot, cfg *config.Config) (IInjector, error) {
	client, err := InitPostgresClient(cfg)
	if err != nil {
		log.Fatal("injector: db client init error")
		return nil, err
	}

  err = CreateTable(log, client)
  if err != nil {
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
