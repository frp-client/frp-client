package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	"github.com/frp-client/frp/client"
	v1 "github.com/frp-client/frp/pkg/config/v1"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	log2 "log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"
)

func (a *App) frpcGracefulClose(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

func (a *App) frpcStartService(cfg *v1.ClientCommonConfig, proxyCfgs []v1.ProxyConfigurer, visitorCfgs []v1.VisitorConfigurer, cfgFile string) error {

	if cfgFile != "" {
		log.Infof("start frpc service for config file [%s]", cfgFile)
		defer log.Infof("frpc service for config file [%s] stopped", cfgFile)
	}
	if a.frpcSvc == nil {
		log2.Println("[准备启动frpcSvc]")
		var err error
		a.frpcSvc, err = client.NewService(client.ServiceOptions{
			Common:         cfg,
			ProxyCfgs:      proxyCfgs,
			VisitorCfgs:    visitorCfgs,
			ConfigFilePath: cfgFile,
		})
		if err != nil {
			return err
		}

		shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
		// Capture the exit signal if we use kcp or quic.
		if shouldGracefulClose {
			go a.frpcGracefulClose(a.frpcSvc)
		}
		return a.frpcSvc.Run(context.Background())
	} else {
		log2.Println("[准备重载frpcSvc]")
		return a.frpcSvc.UpdateAllConfigurer(proxyCfgs, visitorCfgs)
	}
}

func (a *App) FrpcLogin() {
	//util.EmptyOr("", "")
	//a.frpcLogin(utils.ClientId(), utils.ClientId())
}

func (a *App) OnFrpcUpdateConfig(optionalData ...interface{}) {
	log2.Println("[OnFrpcUpdateConfig]", utils.ToJsonString(optionalData))
}

// FrpcStart 启动或者重载proxies
func (a *App) FrpcStart() error {
	var configResp = struct {
		Code uint             `json:"code"`
		Msg  string           `json:"msg"`
		Data model.FrpcConfig `json:"data"`
	}{}
	_, err := utils.HttpJsonGetUnmarshal(
		utils.FormatUrl(apiServer, "/api/frpc/config"),
		a.apiRequestHeaders(),
		&configResp,
	)
	if err != nil {
		log2.Println("[apiConfigError]", err.Error())
		return err
	}
	if configResp.Data.ServerPort <= 0 {
		log2.Println("[apiConfigServerConfigError]", utils.ToJsonString(configResp))
		return err
	}

	a.frpcConfig = &configResp.Data

	var proxiesResp = struct {
		Code uint   `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			List []model.RespUserProxy `json:"list"`
		} `json:"data"`
	}{}
	_, err = utils.HttpJsonGetUnmarshal(
		utils.FormatUrl(apiServer, "/api/frpc/proxies?limit=20"),
		a.apiRequestHeaders(),
		&proxiesResp,
	)
	if err != nil {
		log2.Println("[apiConfigError]", err.Error())
		return err
	}

	a.frpcProxyCfgs = a.parseFrpcProxyConfig(&proxiesResp.Data.List)
	if len(*a.frpcProxyCfgs) == 0 {
		return errors.New("没有可用配置数据")
	}
	if err = a.frpcStartService(a.parseFrpcClientConfig(&configResp.Data), *a.frpcProxyCfgs, make([]v1.VisitorConfigurer, 0), ""); err != nil {
		return err
	}

	return nil
}

func (a *App) apiRequestHeaders(headers ...map[string]string) (h map[string]string) {
	h = make(map[string]string)
	if a.frpcUserSession != nil {
		h["X_CLIENT_ID"] = a.frpcUserSession.MachineId
		h["Authorization"] = fmt.Sprintf("Bearer %s", a.frpcUserSession.JwtToken)
	} else {
		h["X_CLIENT_ID"] = a.ClientId()
	}
	if len(headers) > 0 {
		for key, val := range headers[0] {
			h[key] = val
		}
	}
	return h
}

func (a *App) parseFrpcClientConfig(config *model.FrpcConfig) *v1.ClientCommonConfig {
	var cfg = &v1.ClientCommonConfig{}
	cfg.Complete()

	cfg.ServerAddr = config.ServerAddr
	cfg.ServerPort = int(config.ServerPort)
	cfg.LoginFailExit = &config.LoginFailExit

	cfg.AccessToken = a.frpcUserSession.AccessToken

	return cfg
}

func (a *App) parseFrpcProxyConfig(respProxies *[]model.RespUserProxy) *[]v1.ProxyConfigurer {
	var proxyCfgs = make([]v1.ProxyConfigurer, 0)

	for _, proxy := range *respProxies {
		pc := v1.NewProxyConfigurerByType(v1.ProxyType(proxy.ProxyType.String()))
		tmpHostPort := strings.Split(proxy.ProxyLocalAddr, ":")
		if len(tmpHostPort) != 2 {
			continue
		}
		if proxy.Status != 1 {
			continue
		}

		switch tmpC := pc.(type) {
		case *v1.HTTPProxyConfig:
			tmpC.Name = proxy.ProxyExtra.Subdomain
			tmpC.Type = proxy.ProxyType.String()
			//tmpC.Transport.BandwidthLimitMode = "client"

			tmpC.LocalIP = tmpHostPort[0]
			tmpC.LocalPort = utils.StringToInt(tmpHostPort[1])

			tmpC.CustomDomains = make([]string, 0)
			tmpC.CustomDomains = append(tmpC.CustomDomains, fmt.Sprintf("%s.%s", proxy.ProxyExtra.Subdomain, proxy.Domain))

			proxyCfgs = append(proxyCfgs, pc)
		case *v1.HTTPSProxyConfig:
			certFile, keyFile, err := a.parseCertToFile(
				utils.Md5(fmt.Sprintf("%d,%s", proxy.Id, proxy.SubDomain)),
				[]byte(proxy.ProxyExtra.SslCrt),
				[]byte(proxy.ProxyExtra.SslKey),
			)
			if err != nil {
				continue
			}

			// 参考frp实际运行的配置数据结构填充
			tmpC.Name = proxy.ProxyExtra.Subdomain
			tmpC.Type = proxy.ProxyType.String()
			//tmpC.Transport.BandwidthLimitMode = "client"
			//tmpC.LocalIP = tmpHostPort[0]
			//tmpC.LocalPort = utils.StringToInt(tmpHostPort[1])

			tmpC.CustomDomains = make([]string, 0)
			tmpC.CustomDomains = append(tmpC.CustomDomains, fmt.Sprintf("%s.%s", proxy.ProxyExtra.Subdomain, proxy.Domain))

			tmpC.Plugin.Type = "https2http"
			tmpC.Plugin.ClientPluginOptions = &v1.HTTPS2HTTPPluginOptions{
				Type:              "https2http",
				LocalAddr:         proxy.ProxyLocalAddr,
				CrtPath:           certFile,
				KeyPath:           keyFile,
				HostHeaderRewrite: tmpC.LocalIP,
				RequestHeaders: v1.HeaderOperations{
					Set: map[string]string{
						"x-from-where": "frp",
					},
				},
			}

			proxyCfgs = append(proxyCfgs, pc)
		case *v1.TCPProxyConfig:
			tmpC.Name = proxy.ProxyExtra.Subdomain
			//tmpC.Transport.BandwidthLimitMode = "client"
			tmpC.Type = proxy.ProxyType.String()

			tmpC.LocalIP = tmpHostPort[0]
			tmpC.LocalPort = utils.StringToInt(tmpHostPort[1])
			tmpC.RemotePort = proxy.ProxyRemotePort

			proxyCfgs = append(proxyCfgs, tmpC)
		case *v1.UDPProxyConfig:
			tmpC.Name = proxy.ProxyExtra.Subdomain
			tmpC.Type = proxy.ProxyType.String()

			tmpC.LocalIP = tmpHostPort[0]
			tmpC.LocalPort = utils.StringToInt(tmpHostPort[1])
			tmpC.RemotePort = proxy.ProxyRemotePort

			proxyCfgs = append(proxyCfgs, tmpC)
		default:
		}
	}

	return &proxyCfgs
}

func (a *App) parseCertToFile(id string, certFileBuff, keyFileBuff []byte) (certFile, keyFile string, err error) {
	certFile = utils.AppTempFile("certs", fmt.Sprintf("%sC", id))
	keyFile = utils.AppTempFile("certs", fmt.Sprintf("%sK", id))

	var w = true
	// 比较并修改
	if len(certFileBuff) > 0 && bytes.Compare(utils.ReadFileAsByte(certFile), certFileBuff) == 0 {
		w = false
	}
	if w {
		if err = os.WriteFile(certFile, certFileBuff, fs.ModePerm); err != nil {
			return
		}
	}

	w = true
	// 比较并修改
	if len(keyFile) > 0 && bytes.Compare(utils.ReadFileAsByte(keyFile), keyFileBuff) == 0 {
		w = false
	}
	if w {
		if err = os.WriteFile(keyFile, keyFileBuff, fs.ModePerm); err != nil {
			return
		}
	}
	return
}

func (a *App) startWebServer() error {
	defer func() { recover() }()

	var p = a.appConfig.LocalWebServerPort
	l, err := net.Listen(fiber.NetworkTCP4, fmt.Sprintf(":%d", p))
	if err != nil {
		log2.Println("[本地web服务端口被占用]", err.Error())
		return err
	}
	_ = l.Close()

	app := fiber.New()
	app.Static("/", a.appConfig.LocalWebServerPath, fiber.Static{Browse: true})
	err = app.Listen(fmt.Sprintf(":%d", p))
	if err != nil {
		log2.Println("[本地web服务启动失败]", err.Error())
		return err
	}
	return nil
}

func (a *App) startTcpServer() error {
	defer func() { recover() }()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", a.appConfig.LocalTcpServerPort))
	if err != nil {
		log2.Println("[本地tcp服务端口被占用]", err.Error())
		return err
	}
	var buf []byte
	for {
		accept, err := ln.Accept()
		if err != nil {
			log2.Println("[本地tcp服务Accept]", err.Error())
			continue
		}
		buf = make([]byte, 1024)
		_, err = accept.Read(buf)
		if err != nil {
			log2.Println("[本地tcp服务Read]", err.Error())
			continue
		}
		log2.Println("[本地tcp服务Read]", accept.RemoteAddr().String(), bytes.TrimSpace(buf))
		_, err = accept.Write([]byte(a.appConfig.LocalTcpServerResponse))
		if err != nil {
			log2.Println("[本地tcp服务Write]", err.Error())
			continue
		}
	}
}

func (a *App) startUdpServer() error {
	defer func() { recover() }()
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: a.appConfig.LocalUdpServerPort})
	if err != nil {
		log2.Println("[本地tcp服务启动失败]", err.Error())
		return err
	}
	var buf []byte
	for {
		buf = make([]byte, 1024)
		_, addr, err := ln.ReadFromUDP(buf)
		if err != nil {
			log2.Println("[本地udp服务Read]", err.Error())
			continue
		}
		log2.Println("[本地udp服务Read]", addr.String(), bytes.TrimSpace(buf))
		_, err = ln.WriteToUDP([]byte(a.appConfig.LocalUdpServerResponse), addr)
		if err != nil {
			log2.Println("[本地udp服务Write]", err.Error())
			continue
		}
	}
}
