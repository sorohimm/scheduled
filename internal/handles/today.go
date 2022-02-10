package handles

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func (h *Handles) TodayScheduleInGroup(m *tb.Message) {
	conn, err := h.DBHandler.AcquireConn(context.Background())
	if err != nil {
		h.Log.Info(err.Error())
		_, err = h.Bot.Send(m.Chat, "Бот утонул")
		if err != nil {
			h.Log.Warn(err)
		}
		return
	}
	defer conn.Release()

	group, err := h.getGroup(conn, m.Chat.ID)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		h.Bot.Send(m.Chat, "Группа не установлена")
		return
	case err != nil:
		h.Bot.Send(m.Chat, "Бот утонул")
		return
	}

	id, err := h.getGroupId(group)
	if err != nil {
		_, err = h.Bot.Send(m.Chat, "Нет такой группы, сори(")
		if err != nil {
			h.Log.Warn(err)
		}
		return
	}

	sch, err := h.getSchedule(id)
	if err != nil {
		_, err = h.Bot.Send(m.Chat, "Бот утонул")
		if err != nil {
			h.Log.Warn(err)
		}
		return
	}

	lesns := h.createTodaySchedule(sch.Lessons, int(time.Now().Weekday()), sch.IsNumeratorFirst)
	_, err = h.Bot.Send(m.Chat, lesns)
	if err != nil {
		h.Log.Warn(err)
	}
}

func (h *Handles) getGroup(conn *pgxpool.Conn, chatId int64) (string, error) {
	const GetUserBalanceStatement = `SELECT group_name FROM chats WHERE chat_id = $1;`
	var group string

	err := conn.QueryRow(context.Background(), GetUserBalanceStatement, chatId).Scan(&group)
	if err != nil {
		h.Log.Info(err.Error())
		return "", err
	}

	return group, nil
}
