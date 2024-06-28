package ss

import (
	"sync"
	"time"
)

var (
	tcpOpen    = false
	ssTcpClose = make(chan bool, 1)
	mutex      sync.Mutex
)

func TcpOpen() {
	mutex.Lock()
	defer mutex.Unlock()
	tcpOpen = true
}

func TcpClose() {
	mutex.Lock()
	defer mutex.Unlock()
	if tcpOpen {
		ssTcpClose <- true
		tcpOpen = false
		time.Sleep(time.Second * 2)
	}
}

func TcpStatus() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return tcpOpen
}
