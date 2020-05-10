package define

import "encoding/json"

const (
	RedisRoomConfigKey = "room_config"
)

type RoomConfig struct {
	Master   string `json:"master"`
	Status   int    `json:"status"`
	CurrWord string `json:"currWord"`
	CurrUser string `json:"currUser"`
	Round    int    `json:"round"`
}

func (s RoomConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s RoomConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
