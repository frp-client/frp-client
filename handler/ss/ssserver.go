package ss

import "net"

// https://github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-server/server.go

import (
	"encoding/binary"
	"fmt"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	idType  = 0 // address type index
	idIP0   = 1 // ip address start index
	idDmLen = 1 // domain address length index
	idDm0   = 2 // domain address start index

	typeIPv4 = 1 // type is ipv4 address
	typeDm   = 3 // type is domain address
	typeIPv6 = 4 // type is ipv6 address

	lenIPv4   = net.IPv4len + 2 // ipv4 + 2port
	lenIPv6   = net.IPv6len + 2 // ipv6 + 2port
	lenDmBase = 2               // 1addrLen + 2port, plus addrLen
	// lenHmacSha1 = 10
)

var debug ss.DebugLog
var sanitizeIps bool
var udp bool

func getRequest(conn *ss.Conn) (host string, err error) {
	ss.SetReadTimeout(conn)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 255(max length address) + 2(port) + 10(hmac-sha1)
	buf := make([]byte, 269)
	// read till we get possible domain length field
	if _, err = io.ReadFull(conn, buf[:idType+1]); err != nil {
		return
	}

	var reqStart, reqEnd int
	addrType := buf[idType]
	switch addrType & ss.AddrMask {
	case typeIPv4:
		reqStart, reqEnd = idIP0, idIP0+lenIPv4
	case typeIPv6:
		reqStart, reqEnd = idIP0, idIP0+lenIPv6
	case typeDm:
		if _, err = io.ReadFull(conn, buf[idType+1:idDmLen+1]); err != nil {
			return
		}
		reqStart, reqEnd = idDm0, idDm0+int(buf[idDmLen])+lenDmBase
	default:
		err = fmt.Errorf("addr type %d not supported", addrType&ss.AddrMask)
		return
	}

	if _, err = io.ReadFull(conn, buf[reqStart:reqEnd]); err != nil {
		return
	}

	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch addrType & ss.AddrMask {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+int(buf[idDmLen])])
	}
	// parse port
	port := binary.BigEndian.Uint16(buf[reqEnd-2 : reqEnd])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

const logCntDelta = 100

var connCnt int
var nextLogConnCnt = logCntDelta

func sanitizeAddr(addr net.Addr) string {
	if sanitizeIps {
		return "x.x.x.x:zzzz"
	} else {
		return addr.String()
	}
}

func handleConnection(conn *ss.Conn, port string) {
	var host string

	connCnt++ // this maybe not accurate, but should be enough
	if connCnt-nextLogConnCnt >= 0 {
		// XXX There's no xadd in the atomic package, so it's difficult to log
		// the message only once with low cost. Also note nextLogConnCnt maybe
		// added twice for current peak connection number level.
		log.Printf("Number of client connections reaches %d\n", nextLogConnCnt)
		nextLogConnCnt += logCntDelta
	}

	// function arguments are always evaluated, so surround debug statement
	// with if statement
	if debug {
		debug.Printf("new client %s->%s\n", sanitizeAddr(conn.RemoteAddr()), conn.LocalAddr())
	}
	closed := false
	defer func() {
		if debug {
			debug.Printf("closed pipe %s<->%s\n", sanitizeAddr(conn.RemoteAddr()), host)
		}
		connCnt--
		if !closed {
			conn.Close()
		}
	}()

	host, err := getRequest(conn)
	if err != nil {
		log.Println("error getting request", sanitizeAddr(conn.RemoteAddr()), conn.LocalAddr(), err)
		closed = true
		return
	}
	// ensure the host does not contain some illegal characters, NUL may panic on Win32
	if strings.ContainsRune(host, 0x00) {
		log.Println("invalid domain name.")
		closed = true
		return
	}
	debug.Println("connecting", host)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			log.Println("dial error:", err)
		} else {
			log.Println("error connecting to:", host, err)
		}
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	if debug {
		debug.Printf("piping %s<->%s", sanitizeAddr(conn.RemoteAddr()), host)
	}
	go func() {
		ss.PipeThenClose(conn, remote, func(Traffic int) {
			passwdManager.addTraffic(port, Traffic)
		})
	}()

	ss.PipeThenClose(remote, conn, func(Traffic int) {
		passwdManager.addTraffic(port, Traffic)
	})

	closed = true
	return
}

type PortListener struct {
	password string
	listener net.Listener
}

type UDPListener struct {
	password string
	listener *net.UDPConn
}

type PasswdManager struct {
	sync.Mutex
	portListener map[string]*PortListener
	udpListener  map[string]*UDPListener
	trafficStats map[string]int64
}

func (pm *PasswdManager) add(port, password string, listener net.Listener) {
	pm.Lock()
	pm.portListener[port] = &PortListener{password, listener}
	pm.trafficStats[port] = 0
	pm.Unlock()
}

func (pm *PasswdManager) addUDP(port, password string, listener *net.UDPConn) {
	pm.Lock()
	pm.udpListener[port] = &UDPListener{password, listener}
	pm.Unlock()
}

func (pm *PasswdManager) get(port string) (pl *PortListener, ok bool) {
	pm.Lock()
	pl, ok = pm.portListener[port]
	pm.Unlock()
	return
}

func (pm *PasswdManager) getUDP(port string) (pl *UDPListener, ok bool) {
	pm.Lock()
	pl, ok = pm.udpListener[port]
	pm.Unlock()
	return
}

func (pm *PasswdManager) del(port string) {
	pl, ok := pm.get(port)
	if !ok {
		return
	}
	if udp {
		upl, ok := pm.getUDP(port)
		if !ok {
			return
		}
		upl.listener.Close()
	}
	pl.listener.Close()
	pm.Lock()
	delete(pm.portListener, port)
	delete(pm.trafficStats, port)
	if udp {
		delete(pm.udpListener, port)
	}
	pm.Unlock()
}

func (pm *PasswdManager) addTraffic(port string, n int) {
	pm.Lock()
	pm.trafficStats[port] = pm.trafficStats[port] + int64(n)
	pm.Unlock()
	return
}

func (pm *PasswdManager) getTrafficStats() map[string]int64 {
	pm.Lock()
	copy := make(map[string]int64)
	for k, v := range pm.trafficStats {
		copy[k] = v
	}
	pm.Unlock()
	return copy
}

// Update port password would first close a port and restart listening on that
// port. A different approach would be directly change the password used by
// that port, but that requires **sharing** password between the port listener
// and password manager.
func (pm *PasswdManager) updatePortPasswd(port, password string) {
	pl, ok := pm.get(port)
	if !ok {
		log.Printf("new port %s added\n", port)
	} else {
		if pl.password == password {
			return
		}
		log.Printf("closing port %s to update password\n", port)
		pl.listener.Close()
	}
	// run will add the new port listener to passwdManager.
	// So there maybe concurrent access to passwdManager and we need lock to protect it.
	go run(port, password)
	if udp {
		pl, ok := pm.getUDP(port)
		if !ok {
			log.Printf("new udp port %s added\n", port)
		} else {
			if pl.password == password {
				return
			}
			log.Printf("closing udp port %s to update password\n", port)
			pl.listener.Close()
		}
		go runUDP(port, password)
	}
}

var passwdManager = PasswdManager{
	portListener: map[string]*PortListener{},
	udpListener:  map[string]*UDPListener{},
	trafficStats: map[string]int64{},
}

func RunSsServer(config ss.Config) (ssTcpListener *net.Listener, err error) {
	ssConfig = &config
	ssConfig.Method = strings.ToLower(ssConfig.Method)
	ssTcpListener, err = run(strconv.Itoa(config.LocalPort), config.Password)

	return
}

func run(port, password string) (ssTcpListener *net.Listener, err error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("error listening port %v: %v\n", port, err)
		return nil, err
	}

	ssTcpListener = &ln

	passwdManager.add(port, password, ln)
	var cipher *ss.Cipher
	log.Printf("server listening port %v ...\n", port)
	go func() {
		var _break = false
		for {
			select {
			case <-ssTcpListenerStatus.ChClose:
				_break = true
				_ = ln.Close()
				break
			default:
				conn, err := ln.Accept()
				if err != nil {
					// listener maybe closed to update password
					debug.Printf("accept error: %v\n", err)
					return
				}
				// Creating cipher upon first connection.
				if cipher == nil {
					//log.Println("creating cipher for port:", port)
					cipher, err = ss.NewCipher(ssConfig.Method, password)
					if err != nil {
						log.Printf("Error generating cipher for port: %s %v\n", port, err)
						conn.Close()
						continue
					}
				}
				go handleConnection(ss.NewConn(conn, cipher.Copy()), port)
			}
			if _break {
				break
			}
		}
	}()

	return ssTcpListener, nil
}

func runUDP(port, password string) {
	var cipher *ss.Cipher
	port_i, _ := strconv.Atoi(port)
	log.Printf("listening udp port %v\n", port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv6zero,
		Port: port_i,
	})
	passwdManager.addUDP(port, password, conn)
	if err != nil {
		log.Printf("error listening udp port %v: %v\n", port, err)
		return
	}
	defer conn.Close()
	cipher, err = ss.NewCipher(ssConfig.Method, password)
	if err != nil {
		log.Printf("Error generating cipher for udp port: %s %v\n", port, err)
		conn.Close()
	}
	SecurePacketConn := ss.NewSecurePacketConn(conn, cipher.Copy())
	for {
		if err := ss.ReadAndHandleUDPReq(SecurePacketConn, func(traffic int) {
			passwdManager.addTraffic(port, traffic)
		}); err != nil {
			debug.Printf("udp read error: %v\n", err)
			return
		}
	}
}

var ssConfig *ss.Config
