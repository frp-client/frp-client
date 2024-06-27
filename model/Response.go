package model

type Response struct {
	Ok   bool        `json:"ok"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ProxyExtra struct {
	Subdomain string `json:"subdomain"`
	SslCrt    string `json:"ssl_crt"`
	SslKey    string `json:"ssl_key"`
}

type RespUserProxy struct {
	Id              int64      `json:"id"`
	UserId          int64      `json:"user_id"`
	NodePortId      int64      `json:"node_port_id"`
	ProxyType       ProxyType  `json:"proxy_type"`
	ProxyName       string     `json:"proxy_name"`
	ProxyRemotePort int        `json:"proxy_remote_port"`
	ProxyLocalAddr  string     `json:"proxy_local_addr"`
	Domain          string     `json:"domain"`
	SubDomain       string     `json:"sub_domain"`
	ProxyExtra      ProxyExtra `json:"proxy_extra"`
	Status          int        `json:"status"`
	CreatedAt       int64      `json:"created_at"`
	UpdatedAt       int64      `json:"updated_at"`
}
