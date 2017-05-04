package mailman

import (
	"poste/util"
	"poste/data"
	"net"
)

const (
	TcpType data.ServerType = "tcp"
	WsType data.ServerType = "ws"
)

type Mailman struct {
	ServerType data.ServerType `json:"type"`
	Host       string `json:"host"`
	Port       int `json:"port"`
}

var M = &Mailman{}

func (m *Mailman)Addr() string {
	return util.ToAddr(m.Host, m.Port)
}

func Serve(host string, port int, serverType data.ServerType) {
	if serverType == WsType {
		M.ServerType = serverType
		serveWs(host, port)
	}
}

func beforeServe(addr *net.TCPAddr) {
	M.Host = addr.IP.String()
	M.Port = addr.Port
}