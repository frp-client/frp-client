package main

import (
	"context"
	"fmt"
	"github.com/energye/systray"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	v1 "github.com/frp-client/frp/pkg/config/v1"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"path"
)

// App struct
type App struct {
	ctx             context.Context
	frpcUsersession *model.UserSession
	frpcConfig      *model.FrpcConfig
	frpcProxyCfgs   *[]v1.ProxyConfigurer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go systray.Run(a.systemTray, func() {})

	runtime.EventsOn(ctx, "onAppMounted", a.onAppMounted)

	go func() {
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
		runtime.EventsOn(ctx, "onFrpcUpdateConfig", a.OnFrpcUpdateConfig)

		err = a.OnFrpcNewConfig()
		if err != nil {
			a.WindowMessage(fmt.Sprintf("客户端启动失败：%s", err.Error()), "提示")
			return
		}
	}()
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
	runtime.EventsEmit(a.ctx, "onStartUpEvent", model.Map{
		"baseURL":     baseURL,
		"clientId":    a.ClientId(),
		"_from":       "domReady",
		"userId":      session.UserId,
		"username":    session.Username,
		"machineId":   session.MachineId,
		"updatedAt":   session.UpdatedAt,
		"jwtToken":    session.JwtToken,
		"accessToken": session.AccessToken,
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
