package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"schbot/internal/config"
	"schbot/internal/infrastructure"
	"time"
)

var (
	cfg *config.Config
	ctx context.Context
	b   *tb.Bot
	log *zap.SugaredLogger
)

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("error loading logger: %s", err)
		os.Exit(1)
		return
	}

	log = logger.Sugar()

	cfg = config.New()
	if err != nil {
		log.Fatalf("config init error: %s", err)
	}
	log.Infof("Config loaded:\n%+v", cfg)

	b, err = tb.NewBot(tb.Settings{
		URL: "https://api.telegram.org",

		Token:  cfg.TgToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Microsecond},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	injector, err := infrastructure.Injector(log, b, cfg)
	if err != nil {
		log.Fatal("main: inject failing")
	}
	hands := injector.InjectHandles()

	b.Handle("/sh", hands.GetSchedule)

	b.Start()
}
