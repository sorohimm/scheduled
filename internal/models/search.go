package models

type SearchGroupReq struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type SearchGroupResp struct {
	Items []struct {
		Uuid       string `json:"uuid"`
		Type       string `json:"type"`
		Caption    string `json:"caption"`
		Additional string `json:"additional"`
	} `json:"items"`
	Total int `json:"total"`
}
