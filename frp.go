package main

import (
	"context"
	"encoding/json"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	"github.com/frp-client/frp/client"
	v1 "github.com/frp-client/frp/pkg/config/v1"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"
)

func (a *App) FrpcStart() (resp model.Response) {
	//return startService(cfg, proxyCfgs, visitorCfgs, cfgFilePath)
	buf, err := utils.HttpJsonGet("http://127.0.0.1:3000/api/frpc/config")
	if err != nil {
		resp = model.Response{
			Ok:  false,
			Msg: err.Error(),
		}
		return
	}

	type Resp struct {
		Code uint             `json:"code"`
		Msg  string           `json:"msg"`
		Data model.FrpcConfig `json:"data"`
	}
	var apiResp Resp
	if err = json.Unmarshal(buf, &apiResp); err != nil {
		resp = model.Response{
			Ok:  false,
			Msg: err.Error(),
		}
		return
	}
	if apiResp.Data.ServerPort <= 0 {
		resp = model.Response{
			Ok:   false,
			Msg:  "服务器端口配置错误",
			Data: apiResp,
		}
		return
	}

	var cfg = &v1.ClientCommonConfig{}
	cfg.Complete()

	cfg.ServerAddr = apiResp.Data.ServerAddr
	cfg.ServerPort = int(apiResp.Data.ServerPort)
	cfg.LoginFailExit = &apiResp.Data.LoginFailExit

	cfg.AccessToken = "ed30dfa6e2c4070c7503722c73931ad8"

	var proxyCfgs = make([]v1.ProxyConfigurer, 0)

	{
		pc := v1.NewProxyConfigurerByType(v1.ProxyType(v1.ProxyTypeHTTP))

		switch tmpC := pc.(type) {
		case *v1.HTTPProxyConfig:
			tmpC.Name = "tmpVhostName"
			tmpC.Transport.BandwidthLimitMode = "client"

			tmpC.LocalIP = "127.0.0.1"
			tmpC.LocalPort = 8000

			tmpC.CustomDomains = make([]string, 0)
			tmpC.CustomDomains = append(tmpC.CustomDomains, "www05.frp.local.me")

			proxyCfgs = append(proxyCfgs, pc)
		case *v1.HTTPSProxyConfig:
			//certFile, keyFile, _ := parseCertToFile(vhost.Id, []byte(vhost.CrtPath), []byte(vhost.KeyPath))
			//
			//// 参考frp实际运行的配置数据结构填充
			//tmpC.Name = tmpVhostName
			//tmpC.Transport.BandwidthLimitMode = "client"
			//
			//tmpC.LocalIP = host
			//tmpC.LocalPort = utils.StringToInt(port)
			//
			//tmpC.CustomDomains = make([]string, 0)
			//tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)
			//
			//tmpC.Plugin.Type = "https2http"
			//tmpC.Plugin.ClientPluginOptions = &v1.HTTPS2HTTPPluginOptions{
			//	Type:              "https2http",
			//	LocalAddr:         vhost.LocalAddr,
			//	HostHeaderRewrite: tmpC.LocalIP,
			//	RequestHeaders: v1.HeaderOperations{
			//		Set: map[string]string{
			//			"x-from-where": "frp",
			//		},
			//	},
			//	CrtPath: certFile,
			//	KeyPath: keyFile,
			//}
			//
			//proxyCfg = tmpC
		case *v1.TCPProxyConfig:
			//tmpC.Name = tmpVhostName
			//tmpC.Transport.BandwidthLimitMode = "client"
			//
			//tmpC.LocalIP = host
			//tmpC.LocalPort = utils.StringToInt(port)
			//
			//tmpC.RemotePort = vhost.RemotePort
			//
			//proxyCfg = tmpC
		case *v1.TCPMuxProxyConfig:
			//tmpC.Name = tmpVhostName
			//tmpC.Transport.BandwidthLimitMode = "client"
			//
			//tmpC.LocalIP = host
			//tmpC.LocalPort = utils.StringToInt(port)
			//
			//tmpC.CustomDomains = make([]string, 0)
			//tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)
			//
			//tmpC.Multiplexer = "httpconnect"
			//
			//proxyCfg = tmpC
			//
		default:

		}
	}

	if err = a.frpcStartService(cfg, proxyCfgs, make([]v1.VisitorConfigurer, 0), ""); err != nil {
		resp = model.Response{
			Ok:  false,
			Msg: err.Error(),
		}
		return
	}

	resp = model.Response{
		Ok: true,
	}

	return

}

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
	svr, err := client.NewService(client.ServiceOptions{
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
		go a.frpcGracefulClose(svr)
	}
	return svr.Run(context.Background())
}
