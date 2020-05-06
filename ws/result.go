package ws

type Result struct {
	CMD  string      `json:"cmd"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
