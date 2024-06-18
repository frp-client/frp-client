package model

type FrpcConfig struct {
	ServerAddr    string `json:"server_addr"`
	ServerPort    uint   `json:"server_port"`
	LoginFailExit bool   `json:"login_fail_exit"`
}
