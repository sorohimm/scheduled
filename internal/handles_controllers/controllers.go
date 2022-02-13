package handles_controllers

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"schbot/internal/config"
	"schbot/internal/errors"
	"schbot/internal/interfaces"
	"strconv"
	"strings"
)

type HandlesController struct {
	Log            *zap.SugaredLogger
	Config         *config.Config
	Bot            *tb.Bot
	HandlesService interfaces.IHandlesService
}

//resolveGroup provide group name
func (c *HandlesController) resolveGroup(msg string) string {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	for i := 0; i < 2; i++ {
		scanner.Scan()
	}
	return scanner.Text()
}

//resolveDay provide day of week
func (c *HandlesController) resolveDay(msg string) int64 {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	for i := 0; i < 3; i++ {
		scanner.Scan()
	}
	day, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	return day
}

func (c *HandlesController) GetTodaySchedule(m *tb.Message) {
	group := c.resolveGroup(m.Text)
	c.Log.Info(group)
	c.Log.Info(m.Chat.ID)

	scheduleMsg, err := c.HandlesService.GetTodaySchedule(group)
	if err != nil {
		_, err = c.Bot.Send(m.Chat, "Бот утонул")
		if err != nil {
			c.Log.Warn(err)
		}
		return
	}

	switch err {
	case errors.ErrServerConfused:
		c.Bot.Send(m.Chat, fmt.Sprintf("Сервис недоступен, попробуйте еще раз :("))
	case errors.ErrInvalidGroup:
		c.Bot.Send(m.Chat, fmt.Sprintf("Группа не найдена: %s", group))
	case nil:
		c.Bot.Send(m.Chat, fmt.Sprintf("Текущая группа: %s", scheduleMsg))
	default:
		c.Bot.Send(m.Chat, fmt.Sprintf("Бот утонул"))
	}

	_, err = c.Bot.Send(m.Chat, scheduleMsg)
}

func (c *HandlesController) GetTodayScheduleInChat(m *tb.Message) {
	scheduleMsg, err := c.HandlesService.GetTodayScheduleInChat(m.Chat.ID)
	switch err {
	case errors.ErrServerConfused:
		c.Bot.Send(m.Chat, fmt.Sprintf("Сервис недоступен, попробуйте еще раз :("))
	case errors.ErrGroupNotSet:
		c.Bot.Send(m.Chat, "Группа не установлена")
	case nil:
		c.Bot.Send(m.Chat, scheduleMsg)
	default:
		c.Bot.Send(m.Chat, fmt.Sprintf("Бот утонул"))
	}
}

func (c *HandlesController) GetTomorrowScheduleInChat(m *tb.Message) {
	scheduleMsg, err := c.HandlesService.GetTomorrowScheduleInChat(m.Chat.ID)
	switch err {
	case errors.ErrServerConfused:
		c.Bot.Send(m.Chat, fmt.Sprintf("Сервис недоступен, попробуйте еще раз :("))
	case errors.ErrGroupNotSet:
		c.Bot.Send(m.Chat, "Группа не установлена")
	case nil:
		c.Bot.Send(m.Chat, scheduleMsg)
	default:
		c.Bot.Send(m.Chat, fmt.Sprintf("Бот утонул"))
	}
}

func (c *HandlesController) GetDailySchedule(m *tb.Message) {
	group := c.resolveGroup(m.Text)
	c.Log.Info(group)
	c.Log.Info(m.Chat.ID)

	day := c.resolveDay(m.Text)
	if day > 7 || day < 1 {
		c.Log.Info(day)
		_, err := c.Bot.Send(m.Chat, fmt.Sprintf("У нас что, %d дней в неделе?", day))
		if err != nil {
			c.Log.Warn(err)
		}
		return
	}

	res, err := c.HandlesService.GetDailySchedule(group, day)
	switch err {
	case errors.ErrServerConfused:
		c.Bot.Send(m.Chat, fmt.Sprintf("Сервис недоступен, попробуйте еще раз :("))
	case errors.ErrInvalidGroup:
		c.Bot.Send(m.Chat, fmt.Sprintf("Группа не найдена: %s", group))
	case nil:
		c.Bot.Send(m.Chat, fmt.Sprintf("Текущая группа: %s", res))
	default:
		c.Bot.Send(m.Chat, fmt.Sprintf("Бот утонул"))
	}

}

func (c *HandlesController) SetChatGroup(m *tb.Message) {
	group := c.resolveGroup(m.Text)
	if group == "" {
		return
	}
	c.Log.Info(group)
	c.Log.Info(m.Chat.ID)

	res, err := c.HandlesService.SetGroupChat(m.Chat.ID, group)
	switch err {
	case errors.ErrServerConfused:
		c.Bot.Send(m.Chat, fmt.Sprintf("Сервис недоступен, попробуйте еще раз :("))
	case errors.ErrInvalidGroup:
		c.Bot.Send(m.Chat, fmt.Sprintf("Группа не найдена: %s", group))
	case nil:
		c.Bot.Send(m.Chat, fmt.Sprintf("Текущая группа: %s", res))
	default:
		c.Bot.Send(m.Chat, fmt.Sprintf("Бот утонул"))
	}
}

const title = `❗️Установить группу в чате: 
«/setg 
   РЛ2-42». Группу писать с новой строки.
❗️Расписание на сегодня: 
«/today». Работает только после установки группы в чате!!
❗️Посмотреть расписание любой группы на сегодня: 
«/tsh
РЛ2-42». Группу писать с новой строки.
❗️Посмотреть расписание любой группы в конкретный день: 
«/sh
РЛ2-42
4». Число - день недели по порядку(от 1 до 7). Группу и день недели писать с новой строки.

♻️По всем вопросам обращаться сюда: @sorohimm`

func (c *HandlesController) Start(m *tb.Message) {
	_, err := c.Bot.Send(m.Chat, title)
	if err != nil {
		c.Log.Warn(err.Error())
	}
}
