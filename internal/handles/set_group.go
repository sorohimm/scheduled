package handles

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

	res, err := h.updateChatGroup(conn, m.Chat.ID, group)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		res, err = h.setChatGroup(conn, m.Chat.ID, group)
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

func (h *Handles) updateChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const UpdateAccountStatement = `UPDATE chats SET group_name = $2 WHERE chat_id = $1
								    RETURNING "group_name";`
	var gr string

	err := conn.QueryRow(context.Background(), UpdateAccountStatement, chatId, chatGroup).Scan(&gr)
	if err != nil {
		h.Log.Info(err.Error())
		return "", err
	}

	return gr, nil
}

func (h *Handles) setChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const CreateUserStatement = `INSERT INTO chats (chat_id, group_name) VALUES ($1, $2) 
								 RETURNING "group_name";`

	var gr string
	err := conn.QueryRow(context.Background(), CreateUserStatement, chatId, chatGroup).Scan(&gr)

	if err != nil {
		h.Log.Info(err.Error())
		return "", err
	}

	return gr, nil
}
