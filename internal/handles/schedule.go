package handles

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/interfaces"
	"schbot/internal/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Handles struct {
	Log       *zap.SugaredLogger
	Config    *config.Config
	Bot       *tb.Bot
	Client    *http.Client
	DBHandler interfaces.IDBHandler
}

//resolveGroup provide group name
func (h *Handles) resolveGroup(msg string) string {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	for i := 0; i < 2; i++ {
		scanner.Scan()
	}
	return scanner.Text()
}

//resolveDay provide day of week
func (h *Handles) resolveDay(msg string) int64 {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	for i := 0; i < 3; i++ {
		scanner.Scan()
	}
	day, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	return day
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

	lesns := h.createTodaySchedule(sch.Lessons, int(time.Now().Weekday()), sch.IsNumeratorFirst)
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

	lesns := h.createDailySchedule(sch.Lessons, int(day))
	_, err = h.Bot.Send(m.Chat, lesns)
	if err != nil {
		h.Log.Warn(err)
	}
}

func (h *Handles) createTodaySchedule(lessons []models.Lesson, day int, is_numerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessonsstr string

	count := 1
	if is_numerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
				count = count + 1
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
				count = count + 1
			}
		}
	}

	return lessonsstr
}

func (h *Handles) createDailySchedule(lessons []models.Lesson, day int) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessons_even = "\U0001F976Числитель:\n"
	var lessons_odd = "\n\U0001F975Знаменатель:\n"

	count_e := 1
	count_o := 1

	for _, lesson := range lessons {
		if lesson.Day == day && lesson.IsNumerator {
			lessons_even += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count_e, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
			count_e = count_e + 1
		}
		if lesson.Day == day && !lesson.IsNumerator {
			lessons_odd += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count_o, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
			count_o = count_o + 1
		}
	}

	return lessons_even + lessons_odd
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

//scheduleReq helpful func for getSchedule
func (h *Handles) scheduleReq(req *http.Request) (models.Schedule, error) {
	resp, err := h.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.Schedule{}, err
	}

	sch := models.Schedule{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sch)
	if err != nil {
		h.Log.Info("decode error")
		return models.Schedule{}, err
	}

	return sch, nil
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
		return "", err
	}

	var group_uuid string
	for i := range res.Items {
		if res.Items[i].Caption == groupName {
			group_uuid = res.Items[i].Uuid
		}
	}

	return group_uuid, nil
}

//idReq helpful func for getGroupId
func (h *Handles) idReq(req *http.Request) (models.SearchGroupResp, error) {
	resp, err := h.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.SearchGroupResp{}, err
	}

	res := models.SearchGroupResp{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		h.Log.Info("decode error")
		return models.SearchGroupResp{}, err
	}

	return res, nil
}
