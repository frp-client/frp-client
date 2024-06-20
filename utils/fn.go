package utils

import (
	"github.com/denisbrodbeck/machineid"
)

var (
	_clientId = ""
)

func ClientId() string {
	if _clientId == "" {
		id, _ := machineid.ProtectedID(AppId)
		_clientId = id
	}
	return _clientId
}
