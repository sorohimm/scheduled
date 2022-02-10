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
	"schbot/internal/models"
	"strconv"
	"strings"
)

type Handles struct {
	Log    *zap.SugaredLogger
	Config *config.Config
	Bot    *tb.Bot
	Client *http.Client
}

func (h *Handles) GetSchedule(m *tb.Message) {
	group := h.resolveGroup(m.Text)
	h.Log.Info(group)
	day := h.resolveDay(m.Text)
	if day > 7 || day < 1 {
		h.Log.Info(day)
		h.Bot.Send(m.Chat, fmt.Sprintf("У нас что, %d дней в неделе?", day))
		return
	}

	id, err := h.getGroupId(group)
	h.Log.Info(id)
	if err != nil {
		h.Bot.Send(m.Chat, "Нет такой группы, сори(")
		return
	}

	sch, err := h.getSchedule(id)
	if err != nil {
		h.Bot.Send(m.Chat, "Бот утонул")
		return
	}

	var lessons_e = "\U0001F976Числитель:\n"
	var lessons_o = "\n\U0001F975Знаменатель:\n"

	count_e := 1
	count_o := 1
	for _, lesson := range sch.Lessons {
		if lesson.Day == int(day) && lesson.IsNumerator {
			lessons_e += fmt.Sprintf("%d. %s - %s\n", count_e, lesson.StartAt, lesson.EndAt)
			lessons_e += fmt.Sprintf("\t%s\n", lesson.Name)
			lessons_e += fmt.Sprintf("\tАуд: %s\n", lesson.Cabinet)
			count_e = count_e + 1
		}
		if lesson.Day == int(day) && !lesson.IsNumerator {
			lessons_o += fmt.Sprintf("%d. %s - %s\n", count_o, lesson.StartAt, lesson.EndAt)
			lessons_o += fmt.Sprintf("\t%s\n", lesson.Name)
			lessons_o += fmt.Sprintf("\tАуд: %s\n", lesson.Cabinet)
			count_o = count_o + 1
		}
	}

	h.Bot.Send(m.Chat, lessons_e+lessons_o)

}

func (h *Handles) getSchedule(groupId string) (models.Schedule, error) {
	req, _ := http.NewRequest(http.MethodGet, h.Config.SchPath+groupId, nil)
	h.Log.Info(req.URL)
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

	h.Log.Info("req do ok")

	sch := models.Schedule{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sch)
	if err != nil {
		h.Log.Info("decode error")
		return models.Schedule{}, err
	}

	h.Log.Info("decode ok")
	return sch, nil
}

//resolveGroup provide group name
func (h *Handles) resolveGroup(msg string) string {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	scanner.Scan()
	scanner.Scan()
	return scanner.Text()
}

//resolveDay provide day of week
func (h *Handles) resolveDay(msg string) int64 {
	scanner := bufio.NewScanner(strings.NewReader(msg))
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	day, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	return day
}

//getSchedule provide group uuid
func (h *Handles) getGroupId(groupName string) (string, error) {
	var body = []byte(fmt.Sprintf(`{"query": "%s", "type": "group"}`, groupName))
	req, _ := http.NewRequest(http.MethodPost, h.Config.SgPath, bytes.NewBuffer(body))
	req.Header.Add("x-bb-token", h.Config.BitopToken)

	h.Log.Info(req.URL)
	res, err := h.idReq(req)
	if err != nil {
		h.Log.Info(err)
		return "", err
	}

	if res.Total == 0 {
		h.Log.Info("nil total")
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