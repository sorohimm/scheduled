package handles

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (h *Handles) SetChatGroup(m *tb.Message) {
	group := h.resolveGroup(m.Text)
	if group == "" {
		return
	}
	h.Log.Info(group)
	h.Log.Info(m.Chat.ID)

	_, err := h.getGroupId(group)
	if err != nil {
		h.Bot.Send(m.Chat, fmt.Sprintf("Группа не найдена: %s", group))
		return
	}

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

	res, err := h.DbMagic.UpdateChatGroup(conn, m.Chat.ID, group)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		res, err = h.DbMagic.SetChatGroup(conn, m.Chat.ID, group)
		if err != nil {
			h.Bot.Send(m.Chat, "Бот утонул")
			return
		}
		h.Bot.Send(m.Chat, fmt.Sprintf("Текущая группа: %s", res))
		return
	case err != nil:
		h.Bot.Send(m.Chat, "Бот утонул")
		return
	}

	h.Bot.Send(m.Chat, fmt.Sprintf("Текущая группа: %s", res))
}
