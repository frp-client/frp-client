package model

type Map map[string]any

type LoginResp struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	JwtToken    string `json:"jwt_token"`
}
