package handles_services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"schbot/internal/config"
	er "schbot/internal/errors"
	"schbot/internal/interfaces"
	"time"
)

type HandleService struct {
	Log           *zap.SugaredLogger
	Config        *config.Config
	Client        *http.Client
	DBHandler     interfaces.IDBHandler
	ScheduleMaker interfaces.IScheduleMaker
	HandlesRepo   interfaces.IHandlesRepo
}

func (s *HandleService) GetTodaySchedule(group string) (string, error) {
	groupId, err := s.GetGroupId(group)
	if err != nil {
		return group, er.ErrInvalidGroup
	}

	req, _ := http.NewRequest(http.MethodGet, s.Config.SchPath+groupId, nil)
	req.Header.Add("x-bb-token", s.Config.BitopToken)

	schedule, err := s.HandlesRepo.GetSchedule(req)
	if err != nil {
		return "", er.ErrServerConfused
	}

	lessons := s.ScheduleMaker.CreateTodaySchedule(schedule.Lessons, group, int(time.Now().Weekday()), schedule.IsNumeratorFirst)

	return lessons, nil
}

func (s *HandleService) GetTodayScheduleInChat(chatId int64) (string, error) {
	conn, err := s.DBHandler.AcquireConn(context.Background())
	if err != nil {
		s.Log.Info(err.Error())
		return "", err
	}
	defer conn.Release()

	group, err := s.HandlesRepo.GetGroupName(conn, chatId)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		return "", er.ErrGroupNotSet
	case err != nil:
		return "", er.ErrServerConfused
	}

	groupId, err := s.GetGroupId(group)
	if err != nil {
		return "", er.ErrServerConfused
	}

	req, _ := http.NewRequest(http.MethodGet, s.Config.SchPath+groupId, nil)
	req.Header.Add("x-bb-token", s.Config.BitopToken)

	schedule, err := s.HandlesRepo.GetSchedule(req)
	if err != nil {
		return "", er.ErrServerConfused
	}

	lesns := s.ScheduleMaker.CreateTodaySchedule(schedule.Lessons, group, int(time.Now().Weekday()), schedule.IsNumeratorFirst)

	return lesns, nil
}

func (s *HandleService) GetTomorrowScheduleInChat(chatId int64) (string, error) {
	conn, err := s.DBHandler.AcquireConn(context.Background())
	if err != nil {
		s.Log.Info(err.Error())
		return "", err
	}
	defer conn.Release()

	group, err := s.HandlesRepo.GetGroupName(conn, chatId)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		return "", er.ErrGroupNotSet
	case err != nil:
		return "", er.ErrServerConfused
	}

	groupId, err := s.GetGroupId(group)
	if err != nil {
		return "", er.ErrServerConfused
	}

	req, _ := http.NewRequest(http.MethodGet, s.Config.SchPath+groupId, nil)
	req.Header.Add("x-bb-token", s.Config.BitopToken)

	schedule, err := s.HandlesRepo.GetSchedule(req)
	if err != nil {
		return "", er.ErrServerConfused
	}

	lesns := s.ScheduleMaker.CreateTomorrowSchedule(schedule.Lessons, group, schedule.IsNumeratorFirst)

	return lesns, nil
}

func (s *HandleService) GetDailySchedule(group string, day int64) (string, error) {
	groupId, err := s.GetGroupId(group)
	if err != nil {
		s.Log.Infof("get group id: %s", err.Error())
		return "", er.ErrInvalidGroup
	}

	req, _ := http.NewRequest(http.MethodGet, s.Config.SchPath+groupId, nil)
	req.Header.Add("x-bb-token", s.Config.BitopToken)

	schedule, err := s.HandlesRepo.GetSchedule(req)
	if err != nil {
		return "", er.ErrServerConfused
	}

	lesns := s.ScheduleMaker.CreateDailySchedule(schedule.Lessons, group, int(day))
	if err != nil {
		return "", er.ErrServerConfused
	}
	return lesns, nil
}

func (s *HandleService) SetGroupChat(chatId int64, group string) (string, error) {
	_, err := s.GetGroupId(group)
	if err != nil {
		return "", er.ErrInvalidGroup
	}

	conn, err := s.DBHandler.AcquireConn(context.Background())
	if err != nil {
		s.Log.Info(err.Error())
		return "", er.ErrServerConfused
	}
	defer conn.Release()

	res, err := s.HandlesRepo.UpdateChatGroup(conn, chatId, group)
	switch {
	case errors.Cause(err) == pgx.ErrNoRows:
		res, err = s.HandlesRepo.SetChatGroup(conn, chatId, group)
		if err != nil {
			return "", er.ErrServerConfused
		}
		return res, nil
	case err != nil:
		return "", err
	}

	return res, nil
}

//GetGroupId provide group uuid
func (s *HandleService) GetGroupId(groupName string) (string, error) {
	const groupRequestStatement = `{"query": "%s", "type": "group"}`
	var body = []byte(fmt.Sprintf(groupRequestStatement, groupName))
	req, _ := http.NewRequest(http.MethodPost, s.Config.SgPath, bytes.NewBuffer(body))
	req.Header.Add("x-bb-token", s.Config.BitopToken)

	res, err := s.HandlesRepo.GetGroupUuid(req)
	if err != nil {
		s.Log.Info(err)
		return "", err
	}

	if res.Total == 0 {
		return "", er.ErrNotFound
	}

	var group_uuid string
	for i := range res.Items {
		if res.Items[i].Caption == groupName {
			group_uuid = res.Items[i].Uuid
		}
	}

	return group_uuid, nil
}
