package handles

import (
	"bufio"
	"encoding/json"
	"net/http"
	"schbot/internal/models"
	"strconv"
	"strings"
)

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
