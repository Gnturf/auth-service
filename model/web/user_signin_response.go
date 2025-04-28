package web

type UserSigninResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UUID         string `json:"uuid"`
}
