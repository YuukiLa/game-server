package response

type LoginResp struct {
	Token string `json:"token"`
	UserInfoResp
}

type UserInfoResp struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
