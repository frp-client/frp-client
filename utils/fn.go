package utils

import (
	"github.com/denisbrodbeck/machineid"
)

func ClientId() string {
	id, _ := machineid.ProtectedID(AppId)
	return id
}
