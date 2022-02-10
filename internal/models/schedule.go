package models

type Schedule struct {
	Group            Group    `json:"group"`
	IsNumeratorFirst bool     `json:"is_numerator_first"`
	Lessons          []Lesson `json:"lessons"`
	SemesterEnd      string   `json:"semester_end"`
	SemesterStart    string   `json:"semester_start"`
	Teacher          Teacher  `json:"teacher"`
	Type             string   `json:"type"`
}

type Lesson struct {
	Cabinet     string    `json:"cabinet"`
	Day         int       `json:"day"`
	EndAt       string    `json:"end_at"`
	Groups      []Group   `json:"groups"`
	IsNumerator bool      `json:"is_numerator"`
	Name        string    `json:"name"`
	StartAt     string    `json:"start_at"`
	Teachers    []Teacher `json:"teachers"`
	Type        string    `json:"type"`
	Uuid        string    `json:"uuid"`
}

type Group struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

type Teacher struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

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
