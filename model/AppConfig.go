package model

import "github.com/frp-client/frp-client/utils"

type AppConfig struct {
	ApiServer       string `json:"api_server"`
	LocalServer     bool   `json:"local_server"`
	LocalServerPort int    `json:"local_server_port"`
	LocalServerPath string `json:"local_server_path"`
	Log             bool   `json:"log"`
	LogPath         string `json:"log_path"`
	UpdatedAt       int64  `json:"updated_at"`
}

func NewAppConfig(appConfig AppConfig) *AppConfig {
	return &AppConfig{
		ApiServer:       appConfig.ApiServer,
		LocalServer:     appConfig.LocalServer,
		LocalServerPort: appConfig.LocalServerPort,
		Log:             appConfig.Log,
		LogPath:         appConfig.LogPath,
		UpdatedAt:       utils.UnixTimestamp(),
	}
}

func (x *AppConfig) SetApiServer(apiServer string) *AppConfig {
	x.ApiServer = apiServer
	return x
}

func (x *AppConfig) SetLocalServer(localServer bool) *AppConfig {
	x.LocalServer = localServer
	return x
}

func (x *AppConfig) SetLocalServerPort(localServerPort int) *AppConfig {
	x.LocalServerPort = localServerPort
	return x
}

func (x *AppConfig) SetLog(log bool) *AppConfig {
	x.Log = log
	return x
}

func (x *AppConfig) SetLogFile(logPath string) *AppConfig {
	x.LogPath = logPath
	return x
}
