package model

import (
	"sync"
	"time"
)

type ListenerStatus struct {
	mutex   sync.Mutex
	IsOpen  bool
	ChClose chan bool
}

func NewListenerStatus() *ListenerStatus {
	return &ListenerStatus{
		mutex:   sync.Mutex{},
		IsOpen:  false,
		ChClose: make(chan bool, 1),
	}
}

func (x *ListenerStatus) Open() {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	x.IsOpen = true
}

func (x *ListenerStatus) Close() {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	if x.IsOpen {
		x.ChClose <- true
		x.IsOpen = false
		time.Sleep(time.Second * 2)
	}
}

func (x *ListenerStatus) Status() bool {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	return x.IsOpen
}
