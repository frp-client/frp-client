package udp

import (
	"bytes"
	"github.com/frp-client/frp-client/model"
	log2 "log"
	"net"
)

var udpListenerStatus = model.NewListenerStatus()

func ListenerOpen() bool {
	return udpListenerStatus.IsOpen
}

func ListenerClose() {
	udpListenerStatus.Close()
}

func StartUdpServer(port int, response string) (conn *net.UDPConn, err error) {
	defer func() { recover() }()
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		log2.Println("[本地tcp服务启动失败]", err.Error())
		return nil, err
	}

	go func() {
		var _break = false

		var buf []byte
		for {
			select {
			case <-udpListenerStatus.ChClose:
				_ = ln.Close()
				_break = true
				break
			default:
				buf = make([]byte, 1024)
				_, addr, err := ln.ReadFromUDP(buf)
				if err != nil {
					log2.Println("[本地udp服务Read]", err.Error())
					continue
				}
				log2.Println("[本地udp服务Read]", addr.String(), bytes.TrimSpace(buf))
				_, err = ln.WriteToUDP([]byte(response), addr)
				if err != nil {
					log2.Println("[本地udp服务Write]", err.Error())
					continue
				}
			}
			if _break {
				break
			}
		}
	}()

	return ln, nil
}
