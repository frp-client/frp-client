package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/energye/systray"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	"github.com/frp-client/frp/client"
	v1 "github.com/frp-client/frp/pkg/config/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx             context.Context
	frpcUserSession *model.UserSession
	frpcConfig      *model.FrpcConfig
	appConfig       *model.AppConfig
	frpcProxyCfgs   *[]v1.ProxyConfigurer
	frpcSvc         *client.Service
	svc             struct {
		webServer   *fiber.App
		tcpServer   *net.Listener
		udpServer   *net.UDPConn
		ssTcpServer *net.Listener
		ssUdpServer *net.PacketConn
	}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.InstanceCheck()

	runtime.EventsOn(ctx, "onAppMounted", a.onAppMounted)
	// 注册事件，并启动frpc
	//runtime.EventsOn(ctx, "onFrpcNewConfig", a.OnFrpcNewConfig)
	runtime.EventsOn(a.ctx, "onFrpcUpdateConfig", a.OnFrpcUpdateConfig)

	go systray.Run(a.systemTray, func() {})
	go a.initApp()
}
func (a *App) systemTray() {
	var p = path.Join(utils.AppPath(), "frontend/src/assets/images/instant_mix_24dp.ico")
	systray.SetIcon(utils.ReadFileAsByte(p))

	systray.AddMenuItem("显示", "Show The Window").Click(func() {
		runtime.Show(a.ctx)
	})
	systray.AddSeparator()
	systray.AddMenuItem("退出", "Quit The Program").Click(func() {
		runtime.Quit(a.ctx)
	})
	systray.SetOnDClick(func(menu systray.IMenu) {
		runtime.Show(a.ctx)
	})

	systray.SetOnRClick(func(menu systray.IMenu) { _ = menu.ShowMenu() })
}

func (a *App) InstanceCheck() bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", appInstanceCheckPort))
	if err != nil {
		log.Println("[程序已启动或端口被占用]", err.Error())
		a.WindowMessage("程序已启动或端口被占用", "错误", runtime.ErrorDialog)
		a.Quit()
	}
	go func() { _, _ = ln.Accept() }()
	return true
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Quit() {
	log.Println("[Quit]")
	runtime.Quit(a.ctx)
}

func (a *App) Hidden() {
	runtime.Hide(a.ctx)
}

func (a *App) OpenUrl(openUrl string) {
	utils.OpenUrl(openUrl)
}

func (a *App) ClientId() string {
	return utils.ClientId()
}

func (a *App) onDomReady(ctx context.Context) {
	a.onAppMounted()
}

func (a *App) onAppMounted(optionalData ...interface{}) {
	session, _ := a.checkUserSession()
	config, _ := a.checkAppConfig()
	runtime.EventsEmit(a.ctx, "onStartUpEvent", model.Map{
		"apiServer": apiServer,
		"clientId":  a.ClientId(),
		"_from":     "domReady",
		"session":   session,
		"config":    config,
	})
}

func (a *App) WindowMessage(msg, title string, dialogType ...runtime.DialogType) {
	if len(dialogType) > 0 {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    dialogType[0],
			Title:   title,
			Message: msg,
		})
	} else {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   title,
			Message: msg,
		})
	}
}

// 检查用户session数据
func (a *App) checkAppConfig() (apConfig model.AppConfig, err error) {
	buf, err := utils.AesDecrypt(a.aesKey(), utils.ReadFileAsString(a.getTempConfigFile()))
	if err != nil {
		return apConfig, err
	}
	if err = json.Unmarshal(buf, &apConfig); err != nil {
		return apConfig, err
	}
	return apConfig, nil
}

func (a *App) initAppConfig() (model.AppConfig, error) {
	config, err := a.checkAppConfig()
	if err != nil {
		//初始化
		config = model.AppConfig{
			ApiServer:              apiServer,
			LocalWebServer:         true,
			LocalWebServerPort:     localWebServerPort,
			LocalWebServerPath:     utils.AppPath(),
			LocalTcpServer:         false,
			LocalTcpServerPort:     localTcpServerPort,
			LocalTcpServerResponse: "tcp ok",
			LocalUdpServer:         false,
			LocalUdpServerPort:     localUdpServerPort,
			LocalUdpServerResponse: "udp ok",
			LocalSsServer:          false,
			LocalSsPort:            localSsServerPort,
			LocalSsCipher:          "AES-256-GCM",
			LocalSsPassword:        a.ClientId()[:8],
			Log:                    false,
			LogPath:                utils.AppPath(),
			UpdatedAt:              utils.UnixTimestamp(),
		}

		//utils.AesDecrypt
		encrypt, err := utils.AesEncrypt(a.aesKey(), utils.ToJsonString(config))
		if err != nil {
			return config, err
		}
		err = utils.SaveFileAsString(a.getTempConfigFile(), encrypt)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

func (a *App) initApp() {
	//initAppConfig
	config, err := a.initAppConfig()
	if err != nil {
		return
	}
	a.appConfig = &config

	// 是否启动本地web服务
	a.reloadSettingConfig()

	// 准备初始化并启动frpc
	session, err := a.apiClientLogin()
	if err != nil {
		log.Println("[客户端登陆失败]", err.Error())
		//runtime.Quit(ctx)
		a.WindowMessage(fmt.Sprintf("客户端授权失败：%s", err.Error()), "提示")
		a.Quit()
	}
	a.frpcUserSession = &session
	err = a.FrpcStart()
	if err != nil {
		a.WindowMessage(fmt.Sprintf("客户端启动失败：%s", err.Error()), "提示")
		a.Quit()
	}

}

func (a *App) AppConfigUpdate(appConfig model.AppConfig) error {
	if appConfig.LocalWebServer {
		if appConfig.LocalWebServerPort <= 0 || appConfig.LocalWebServerPort >= 65535 {
			return errors.New("WEB端口配置不合法")
		}
	}
	if appConfig.LocalTcpServer {
		if appConfig.LocalTcpServerPort <= 0 || appConfig.LocalTcpServerPort >= 65535 {
			return errors.New("TCP端口配置不合法")
		}
	}
	if appConfig.LocalUdpServer {
		if appConfig.LocalUdpServerPort <= 0 || appConfig.LocalUdpServerPort >= 65535 {
			return errors.New("UDP端口配置不合法")
		}
	}
	if appConfig.LocalSsServer {
		if appConfig.LocalSsPort <= 0 || appConfig.LocalSsPort >= 65535 {
			return errors.New("SSServer端口配置不合法")
		}
	}

	stat, err := os.Stat(appConfig.LocalWebServerPath)
	if err != nil {
		return errors.New(fmt.Sprintf("WEB根目录异常：%s", err.Error()))
	}
	if !stat.IsDir() {
		return errors.New("WEB根目录异常")
	}

	{
		//api/version
		var versionResp = struct {
			Code uint   `json:"code"`
			Msg  string `json:"msg"`
			Data struct {
				Version string `json:"version"`
			} `json:"data"`
		}{}
		_, err := utils.HttpJsonGetUnmarshal(
			utils.FormatUrl(appConfig.ApiServer, "/api/version"),
			a.apiRequestHeaders(),
			&versionResp,
		)
		if err != nil {
			return errors.New(fmt.Sprintf("API服务器异常：%s", err.Error()))
		}
		if versionResp.Code != 200 {
			return errors.New(fmt.Sprintf("API服务器异常：%s", utils.ToJsonString(versionResp)))
		}
	}

	config, err := a.checkAppConfig()
	if err != nil {
		config = model.AppConfig{
			ApiServer:          appConfig.ApiServer,
			LocalWebServer:     appConfig.LocalWebServer,
			LocalWebServerPort: appConfig.LocalWebServerPort,
			LocalWebServerPath: appConfig.LocalWebServerPath,

			LocalTcpServer:         appConfig.LocalSsServer,
			LocalTcpServerPort:     appConfig.LocalTcpServerPort,
			LocalTcpServerResponse: appConfig.LocalTcpServerResponse,
			LocalUdpServer:         appConfig.LocalUdpServer,
			LocalUdpServerPort:     appConfig.LocalSsPort,
			LocalUdpServerResponse: appConfig.LocalUdpServerResponse,
			LocalSsServer:          appConfig.LocalSsServer,
			LocalSsPort:            appConfig.LocalSsPort,
			LocalSsCipher:          appConfig.LocalSsCipher,
			LocalSsPassword:        appConfig.LocalSsPassword,
			Log:                    appConfig.Log,
			LogPath:                utils.AppPath(),
			UpdatedAt:              utils.UnixTimestamp(),
		}
	} else {
		config.ApiServer = appConfig.ApiServer
		config.LocalWebServer = appConfig.LocalWebServer
		config.LocalWebServerPort = appConfig.LocalWebServerPort
		config.LocalWebServerPath = appConfig.LocalWebServerPath
		config.LocalTcpServer = appConfig.LocalTcpServer
		config.LocalTcpServerPort = appConfig.LocalTcpServerPort
		config.LocalTcpServerResponse = appConfig.LocalTcpServerResponse
		config.LocalUdpServer = appConfig.LocalUdpServer
		config.LocalUdpServerPort = appConfig.LocalUdpServerPort
		config.LocalUdpServerResponse = appConfig.LocalUdpServerResponse
		config.LocalSsServer = appConfig.LocalSsServer
		config.LocalSsPort = appConfig.LocalSsPort
		config.LocalSsCipher = appConfig.LocalSsCipher
		config.LocalSsPassword = appConfig.LocalSsPassword
		config.Log = appConfig.Log
		config.UpdatedAt = utils.UnixTimestamp()
	}

	//utils.AesDecrypt
	encrypt, err := utils.AesEncrypt(a.aesKey(), utils.ToJsonString(config))
	if err != nil {
		log.Println("[AppConfig数据编码失败]", err.Error())
		return err
	}
	err = utils.SaveFileAsString(a.getTempConfigFile(), encrypt)
	if err != nil {
		log.Println("[AppConfig数据写文件失败]", err.Error())
		return err
	}

	a.appConfig = &config
	a.reloadSettingConfig()

	return nil
}

func (a *App) reloadSettingConfig() {
	// 是否启动本地web服务
	if a.appConfig.LocalWebServer == true {
		go func() { _ = a.RpcStartWebServer() }()
	} else {
		go func() { _ = a.RpcStopWebServer() }()
	}
	if a.appConfig.LocalTcpServer == true {
		go func() { _ = a.RpcStartTcpServer() }()
	} else {
		go func() { _ = a.RpcStopTcpServer() }()
	}
	if a.appConfig.LocalUdpServer == true {
		go func() { _ = a.RpcStartUdpServer() }()
	} else {
		go func() { _ = a.RpcStopUdpServer() }()
	}
	if a.appConfig.LocalSsServer == true {
		go func() { _ = a.RpcStartSsServer() }()
	} else {
		go func() { _ = a.RpcStopSsServer() }()
	}
	a.setupLog()
}

func (a *App) setupLog() {
	if a.appConfig.Log == true {
		var logFile = filepath.Join(a.appConfig.LogPath, fmt.Sprintf("frp-client-%s.log", time.Now().Format("200612")))
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Println("[日志文件打开失败]", logFile, err.Error())
		}
	} else {
		log.SetOutput(io.Discard)
	}
}

func (a *App) AppConfig() model.AppConfig {
	config, _ := a.checkAppConfig()
	return config
}

func (a *App) AppConfigReset() error {
	if err := os.Remove(a.getTempConfigFile()); err != nil {
		return err
	}
	_, err := a.initAppConfig()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) AppLogs(startLine int) model.Map {
	var logFile = filepath.Join(a.appConfig.LogPath, fmt.Sprintf("frp-client-%s.log", time.Now().Format("200612")))
	file, err := os.Open(logFile)
	if err != nil {
		log.Println("[日志文件打开失败]", logFile, err.Error())
		return nil
	}
	defer func() { _ = file.Close() }()

	var lastXLine = 100
	//var lines []string
	if startLine <= 0 {
		// 获取文件总行数，并返回最后100行
		var totalScanner = bufio.NewScanner(file)
		var totalLine = 0
		for totalScanner.Scan() {
			totalLine++
		}
		if err := totalScanner.Err(); err != nil {
			return nil
		}
		if totalLine >= lastXLine*2 {
			startLine = totalLine - lastXLine
		} else {
			startLine = 0
		}
	}

	var readFromLine = func(filePath string, lineNum int) (lastLine int, lines []string, err error) {
		file, err := os.Open(filePath)
		if err != nil {
			return 0, nil, err
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)

		var maxScanLine = 50
		for scanner.Scan() {
			lastLine++
			if lastLine > lineNum {
				lines = append(lines, scanner.Text())
			}
			if lastLine > lineNum+maxScanLine {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return 0, nil, err
		}
		return
	}

	lastLine, lines, err := readFromLine(logFile, startLine)
	if err != nil {
		return nil
	}

	return model.Map{
		"start_line": startLine,
		"last_line":  lastLine,
		"logs":       lines,
	}
}
