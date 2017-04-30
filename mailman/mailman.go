package mailman

import "poste/util"

type ServerType string

const (
	TcpType ServerType = "tcp"
	WsType ServerType = "ws"
)

type Mailman struct {
	ServerType ServerType `json:"type"`
	Host       string `json:"host"`
	Port       int `json:"port"`
}

var M = &Mailman{}

func (m *Mailman)Addr() string {
	return util.ToAddr(m.Host, m.Port)
}

func Serve(host string, port int, serverType ServerType) {
	if serverType == WsType {
		M.ServerType = serverType
		serveWs(host, port)
	}
}

