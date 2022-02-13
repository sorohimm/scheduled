package handles_repos

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"schbot/internal/config"
	"schbot/internal/models"
)

type HandlesRepo struct {
	Log    *zap.SugaredLogger
	Config *config.Config
	Client *http.Client
}

func (r *HandlesRepo) UpdateChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const UpdateAccountStatement = `UPDATE chats SET group_name = $2 WHERE chat_id = $1
								    RETURNING "group_name";`
	var gr string

	err := conn.QueryRow(context.Background(), UpdateAccountStatement, chatId, chatGroup).Scan(&gr)
	if err != nil {
		r.Log.Info(err.Error())
		return "", err
	}

	return gr, nil
}

func (r *HandlesRepo) SetChatGroup(conn *pgxpool.Conn, chatId int64, chatGroup string) (string, error) {
	const CreateUserStatement = `INSERT INTO chats (chat_id, group_name) VALUES ($1, $2) 
								 RETURNING "group_name";`

	var gr string
	err := conn.QueryRow(context.Background(), CreateUserStatement, chatId, chatGroup).Scan(&gr)

	if err != nil {
		r.Log.Warn("db query error: %s", err.Error())
		return "", err
	}

	return gr, nil
}

func (r *HandlesRepo) GetGroupName(conn *pgxpool.Conn, chatId int64) (string, error) {
	const GetUserBalanceStatement = `SELECT group_name FROM chats WHERE chat_id = $1;`
	var group string

	err := conn.QueryRow(context.Background(), GetUserBalanceStatement, chatId).Scan(&group)
	if err != nil {
		r.Log.Warn("db query error: %s", err.Error())
		return "", err
	}

	return group, nil
}

//GetGroupUuid helpful func for getGroupId
func (r *HandlesRepo) GetGroupUuid(req *http.Request) (models.SearchGroupResp, error) {
	resp, err := r.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		r.Log.Warn("request error: %s", err.Error())
		return models.SearchGroupResp{}, err
	}

	res := models.SearchGroupResp{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		r.Log.Infof("decode error: %s", err.Error())
		return models.SearchGroupResp{}, err
	}

	return res, nil
}

//GetSchedule helpful func for getSchedule
func (r *HandlesRepo) GetSchedule(req *http.Request) (models.Schedule, error) {
	resp, err := r.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		r.Log.Warn("request error: %s", err.Error())
		return models.Schedule{}, err
	}

	sch := models.Schedule{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sch)
	if err != nil {
		r.Log.Info("decode error")
		return models.Schedule{}, err
	}

	return sch, nil
}
