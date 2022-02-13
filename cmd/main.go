package main

import (
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
		log.Fatalf("configs init error: %s", err)
	}
	log.Infof("Config loaded:\n%+v", cfg)

	b, err = tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
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
	log.Info("inject ok")
	hands := injector.InjectHandles()

	b.Handle("/start", hands.Start)
	b.Handle("/help", hands.Start)
	b.Handle("/sh", hands.GetDailySchedule)
	b.Handle("/tsh", hands.GetTodaySchedule)
	b.Handle("/setg", hands.SetChatGroup)
	b.Handle("/today", hands.GetTodayScheduleInChat)
	b.Handle("/tomorrow", hands.GetTomorrowScheduleInChat)

	b.Start()
}
