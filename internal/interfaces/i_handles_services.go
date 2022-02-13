package interfaces

type IHandlesService interface {
	GetTodaySchedule(groupId string) (string, error)
	GetGroupId(groupName string) (string, error)
	SetGroupChat(chatId int64, group string) (string, error)
	GetTodayScheduleInChat(chatId int64) (string, error)
	GetTomorrowScheduleInChat(chatId int64) (string, error)
	GetDailySchedule(group string, day int64) (string, error)
}
