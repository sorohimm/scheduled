package interfaces

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type IHandles interface {
	GetSchedule(message *tb.Message)
}
