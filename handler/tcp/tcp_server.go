package tcp

import (
	"bytes"
	"github.com/frp-client/frp-client/model"
	log2 "log"
	"net"
)

var tcpListenerStatus = model.NewListenerStatus()

func ListenerOpen() bool {
	return tcpListenerStatus.IsOpen
}

func ListenerClose() {
	tcpListenerStatus.Close()
}

func StartTcpServer(addr, response string) (listener *net.Listener, err error) {
	defer func() { recover() }()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log2.Println("[本地tcp服务端口被占用]", err.Error())
		return nil, err
	}
	listener = &ln

	go func() {
		var _break = false

		var buf []byte
		for {
			select {
			case <-tcpListenerStatus.ChClose:
				_ = ln.Close()
				_break = true
				break
			default:
				accept, err := ln.Accept()
				if err != nil {
					log2.Println("[本地tcp服务Accept!!!!]", err.Error())
					return
				}
				buf = make([]byte, 1024)
				_, err = accept.Read(buf)
				if err != nil {
					log2.Println("[本地tcp服务Read]", err.Error())
					continue
				}
				log2.Println("[本地tcp服务Read]", accept.RemoteAddr().String(), bytes.TrimSpace(buf))
				_, err = accept.Write([]byte(response))
				if err != nil {
					log2.Println("[本地tcp服务Write]", err.Error())
					continue
				}
			}
			if _break {
				break
			}
		}

	}()

	return listener, nil
}
