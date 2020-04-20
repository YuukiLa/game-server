package request

type UserRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
