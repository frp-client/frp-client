package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/energye/systray"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	"github.com/frp-client/frp/client"
	v1 "github.com/frp-client/frp/pkg/config/v1"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx             context.Context
	frpcUsersession *model.UserSession
	frpcConfig      *model.FrpcConfig
	appConfig       *model.AppConfig
	frpcProxyCfgs   *[]v1.ProxyConfigurer
	frpcSvc         *client.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "onAppMounted", a.onAppMounted)

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
			ApiServer:       apiServer,
			LocalServer:     true,
			LocalServerPort: localWebServerPort,
			LocalServerPath: utils.AppPath(),
			Log:             true,
			LogPath:         utils.AppPath(),
			UpdatedAt:       utils.UnixTimestamp(),
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

	if config.Log == true {
		var logFile = filepath.Join(config.LogPath, fmt.Sprintf("frp-client-%s.log", time.Now().Format("200612")))
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Println("[OpenLogFileError]", err.Error())
			return
		}
		log.SetOutput(file)
		log.Println("[SetLogFile]")
	} else {
		log.SetOutput(io.Discard)
	}

	// 是否启动本地web服务
	if config.LocalServer == true {
		go a.startWebServer()
	}

	// 准备初始化并启动frpc
	session, err := a.apiClientLogin()
	if err != nil {
		log.Println("[客户端登陆失败]", err.Error())
		//runtime.Quit(ctx)
		a.WindowMessage(fmt.Sprintf("客户端授权失败：%s", err.Error()), "提示")
		return
	}
	a.frpcUsersession = &session

	// 注册事件，并启动frpc
	//runtime.EventsOn(ctx, "onFrpcNewConfig", a.OnFrpcNewConfig)
	runtime.EventsOn(a.ctx, "onFrpcUpdateConfig", a.OnFrpcUpdateConfig)

	err = a.FrpcStart()
	if err != nil {
		a.WindowMessage(fmt.Sprintf("客户端启动失败：%s", err.Error()), "提示")
		return
	}

}

func (a *App) AppConfigUpdate(appConfig model.AppConfig) error {
	log.Println("[========>appConfig]", utils.ToJsonString(appConfig))
	if appConfig.LocalServer {
		if appConfig.LocalServerPort <= 0 || appConfig.LocalServerPort >= 65535 {
			return errors.New("端口配置不合法")
		}
	}

	stat, err := os.Stat(appConfig.LocalServerPath)
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
			ApiServer:       appConfig.ApiServer,
			LocalServer:     appConfig.LocalServer,
			LocalServerPort: appConfig.LocalServerPort,
			LocalServerPath: appConfig.LocalServerPath,
			Log:             appConfig.Log,
			LogPath:         utils.AppPath(),
			UpdatedAt:       utils.UnixTimestamp(),
		}
	} else {
		config.ApiServer = appConfig.ApiServer
		config.LocalServer = appConfig.LocalServer
		config.LocalServerPort = appConfig.LocalServerPort
		config.LocalServerPath = appConfig.LocalServerPath
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

	return nil
}
