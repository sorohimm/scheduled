package interfaces

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type IHandles interface {
	Start(m *tb.Message)
	GetDailySchedule(message *tb.Message)
	GetTodaySchedule(m *tb.Message)
	SetChatGroup(m *tb.Message)
	TodayScheduleInGroup(m *tb.Message)
}
