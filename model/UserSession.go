package model

import "github.com/frp-client/frp-client/utils"

type UserSession struct {
	UserId      int64  `json:"user_id"`
	Username    string `json:"username"`
	MachineId   string `json:"machine_id"`
	JwtToken    string `json:"jwt_token"`
	AccessToken string `json:"access_token"`
	UpdatedAt   int64  `json:"updated_at"`
}

func NewUserSession() *UserSession {
	return &UserSession{
		UpdatedAt: utils.UnixTimestamp(),
	}
}

func (x *UserSession) SetUserId(userId int64) *UserSession {
	x.UserId = userId
	return x
}

func (x *UserSession) SetMachineId(machineId string) *UserSession {
	x.MachineId = machineId
	return x
}

func (x *UserSession) SetUsername(username string) *UserSession {
	x.Username = username
	return x
}

func (x *UserSession) SetAccessToken(accessToken string) *UserSession {
	x.AccessToken = accessToken
	return x
}

func (x *UserSession) SetJwtToken(jwtToken string) *UserSession {
	x.JwtToken = jwtToken
	return x
}
