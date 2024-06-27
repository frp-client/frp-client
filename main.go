package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
	apiServer          = "https://tunnel-api.lixiang4u.xyz:3000"
	localWebServerPort = 8090
)

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "frp客户端",
		Width:  1024,
		Height: 668,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		OnDomReady:    app.onDomReady,
		Frameless:     true,
		DisableResize: true,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
