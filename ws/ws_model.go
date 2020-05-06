package ws

type RoomBroadcast struct {
	RoomID   string `json:"roomId"`
	Data     Result `json:"data"`
	ClientID string `json:"clientId"`
}
