package model

import "github.com/frp-client/frp-client/utils"

type LocalStorage struct {
	Id        string `json:"id"`
	MachineId string `json:"machine_id"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		UpdatedAt: utils.UnixTimestamp(),
	}
}

func (x *LocalStorage) SetId(id string) *LocalStorage {
	x.Id = id
	return x
}

func (x *LocalStorage) SetMachineId(machineId string) *LocalStorage {
	x.MachineId = machineId
	return x
}

func (x *LocalStorage) SetUsername(username string) *LocalStorage {
	x.Username = username
	return x
}

func (x *LocalStorage) SetToken(token string) *LocalStorage {
	x.Token = token
	return x
}
