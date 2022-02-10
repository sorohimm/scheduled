package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"schbot/internal/config"
)

type DMagic struct {
	Log    *zap.SugaredLogger
	Config *config.Config
}

func (m *DMagic) UpdateChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const UpdateAccountStatement = `UPDATE chats SET group_name = $2 WHERE chat_id = $1
								    RETURNING "group_name";`
	var gr string

	err := conn.QueryRow(context.Background(), UpdateAccountStatement, chatId, chatGroup).Scan(&gr)
	if err != nil {
		m.Log.Info(err.Error())
		return "", err
	}

	return gr, nil
}

func (m *DMagic) SetChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const CreateUserStatement = `INSERT INTO chats (chat_id, group_name) VALUES ($1, $2) 
								 RETURNING "group_name";`

	var gr string
	err := conn.QueryRow(context.Background(), CreateUserStatement, chatId, chatGroup).Scan(&gr)

	if err != nil {
		m.Log.Info(err.Error())
		return "", err
	}

	return gr, nil
}

func (m *DMagic) GetGroupName(conn *pgxpool.Conn, chatId int64) (string, error) {
	const GetUserBalanceStatement = `SELECT group_name FROM chats WHERE chat_id = $1;`
	var group string

	err := conn.QueryRow(context.Background(), GetUserBalanceStatement, chatId).Scan(&group)
	if err != nil {
		m.Log.Info(err.Error())
		return "", err
	}

	return group, nil
}
