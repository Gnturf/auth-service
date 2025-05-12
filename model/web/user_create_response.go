package web

type UserCreateResponse struct {
	EmailToken string `json:"email_token"`
	UUID       string `json:"uuid"`
}