package handles

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/interfaces"
	"schbot/internal/models"
	"time"
)

type Handles struct {
	Log       *zap.SugaredLogger
	Config    *config.Config
	Bot       *tb.Bot
	Client    *http.Client
	DBHandler interfaces.IDBHandler
	DbMagic   interfaces.IDbMagic
}

func (h *Handles) GetTodaySchedule(m *tb.Message) {
	group := h.resolveGroup(m.Text)
	h.Log.Info(group)
	h.Log.Info(m.Chat.ID)

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

	lesns := h.createTodaySchedule(sch.Lessons, group, int(time.Now().Weekday()), sch.IsNumeratorFirst)
	_, err = h.Bot.Send(m.Chat, lesns)
	if err != nil {
		h.Log.Warn(err)
	}
}

func (h *Handles) GetDailySchedule(m *tb.Message) {
	group := h.resolveGroup(m.Text)
	h.Log.Info(group)
	h.Log.Info(m.Chat.ID)

	day := h.resolveDay(m.Text)
	if day > 7 || day < 1 {
		h.Log.Info(day)
		_, err := h.Bot.Send(m.Chat, fmt.Sprintf("У нас что, %d дней в неделе?", day))
		if err != nil {
			h.Log.Warn(err)
		}
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

	lesns := h.createDailySchedule(sch.Lessons, group, int(day))
	_, err = h.Bot.Send(m.Chat, lesns)
	if err != nil {
		h.Log.Warn(err)
	}
}

func (h *Handles) getSchedule(groupId string) (models.Schedule, error) {
	req, _ := http.NewRequest(http.MethodGet, h.Config.SchPath+groupId, nil)
	req.Header.Add("x-bb-token", h.Config.BitopToken)

	res, err := h.scheduleReq(req)
	if err != nil {
		return models.Schedule{}, err
	}

	return res, nil
}

//getSchedule provide group uuid
func (h *Handles) getGroupId(groupName string) (string, error) {
	var body = []byte(fmt.Sprintf(`{"query": "%s", "type": "group"}`, groupName))
	req, _ := http.NewRequest(http.MethodPost, h.Config.SgPath, bytes.NewBuffer(body))
	req.Header.Add("x-bb-token", h.Config.BitopToken)

	res, err := h.idReq(req)
	if err != nil {
		h.Log.Info(err)
		return "", err
	}

	if res.Total == 0 {
		return "", errors.New("empty result")
	}

	var group_uuid string
	for i := range res.Items {
		if res.Items[i].Caption == groupName {
			group_uuid = res.Items[i].Uuid
		}
	}

	return group_uuid, nil
}
