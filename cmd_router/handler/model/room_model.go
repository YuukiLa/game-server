package model

type RoomID struct {
	RoomID string `json:"roomId"`
}

type Danmu struct {
	Danmu    string `json:"danmu"`
	User     string `json:"user"`
	Avatar   string `json:"avatar"`
	IsAnswer bool   `json:"isAnswer"`
}
