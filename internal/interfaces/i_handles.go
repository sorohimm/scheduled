package interfaces

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type IHandles interface {
	GetDailySchedule(message *tb.Message)
	GetTodaySchedule(m *tb.Message)
}
