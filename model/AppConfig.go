package model

type AppConfig struct {
	ApiServer              string `json:"api_server"`
	LocalWebServer         bool   `json:"local_web_server"`
	LocalWebServerPort     int    `json:"local_web_server_port"`
	LocalWebServerPath     string `json:"local_web_server_path"`
	LocalTcpServer         bool   `json:"local_tcp_server"`
	LocalTcpServerPort     int    `json:"local_tcp_server_port"`
	LocalTcpServerResponse string `json:"local_tcp_server_response"`
	LocalUdpServer         bool   `json:"local_udp_server"`
	LocalUdpServerPort     int    `json:"local_udp_server_port"`
	LocalUdpServerResponse string `json:"local_udp_server_response"`
	Log                    bool   `json:"log"`
	LogPath                string `json:"log_path"`
	UpdatedAt              int64  `json:"updated_at"`
}
