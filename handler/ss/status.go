package ss

import "github.com/frp-client/frp-client/model"

var ssTcpListenerStatus = model.NewListenerStatus()

func TcpListenerOpen() bool {
	return ssTcpListenerStatus.IsOpen
}

func TcpListenerClose() {
	ssTcpListenerStatus.Close()
}
