package main

import (
	"context"
	"fmt"
	"frp-client/utils"
	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"path"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(a.systemTray, func() {})
}
func (a *App) systemTray() {
	var p = path.Join(utils.AppPath(), "frontend/src/assets/images/instant_mix_24dp.ico")
	systray.SetIcon(utils.ReadFile(p))

	systray.AddMenuItem("显示", "Show The Window").Click(func() {
		runtime.Show(a.ctx)
	})
	systray.AddSeparator()
	systray.AddMenuItem("退出", "Quit The Program").Click(func() {
		runtime.Show(a.ctx)
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
