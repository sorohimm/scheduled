package checker

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"schbot/internal/config"
	"schbot/internal/interfaces"
)

type Checker struct {
	Log       *zap.SugaredLogger
	Config    *config.Config
	Bot       *tb.Bot
	DBHandler interfaces.IDBHandler
	DbMagic   interfaces.IDbMagic
}

type Notif struct {
	ChatId    int64
	UpdTime   string
	GroupName string
}

func (c *Checker) Start() {
	notifList := make(chan []Notif)

	go c.pull()
	for {

	}
}

func (c *Checker) processUpdate(updList []Notif) {

}

func (c *Checker) pull(ch chan<- []Notif) {
	for {

	}
}

func (c *Checker) getScheduleNotifList(conn *pgxpool.Conn) ([]Notif, error) {
	const notifListStatement = `SELECT chat_id, group_name, notif_time FROM chats`

	var notifList []Notif
	rows, err := conn.Query(context.Background(), notifListStatement)
	if err != nil {
		c.Log.Info(err.Error())
		return nil, err
	}

	for rows.Next() {
		var newNotif Notif
		err := rows.Scan(&newNotif.ChatId, &newNotif.GroupName, &newNotif.UpdTime)
		if err != nil {
			c.Log.Info(err.Error())
			return nil, err
		}
		notifList = append(notifList, newNotif)
	}

	return notifList, nil
}
