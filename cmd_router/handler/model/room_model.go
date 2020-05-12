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

type OtherPlayer struct {
	CurrUser string `json:"currUser"`
	Round    int    `json:"round"`
}

type CurrPlayer struct {
	CurrWord string `json:"currWord"`
}
