package interfaces

import "github.com/jackc/pgx/v4/pgxpool"

type IDbMagic interface {
	UpdateChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error)
	SetChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error)
	GetGroupName(conn *pgxpool.Conn, chatId int64) (string, error)
}
