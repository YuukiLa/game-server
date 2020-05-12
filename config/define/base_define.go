package define

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/mongo/model"
)

const (
	RedisRoomConfigKey = "room_config"
)

type RoomConfig struct {
	Master   string            `json:"master"`
	Status   int               `json:"status"`
	CurrWord string            `json:"currWord"`
	AllWord  []model.WordModel `json:"allWord"`
	CurrUser string            `json:"currUser"`
	Round    int               `json:"round"`
	AllRound int               `json:"allRound"`
}

func (s RoomConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s RoomConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
