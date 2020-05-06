package response

import (
	"encoding/json"
)

type LoginResp struct {
	Token string `json:"token"`
	UserInfoResp
}

type UserInfoResp struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (s UserInfoResp) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s UserInfoResp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
