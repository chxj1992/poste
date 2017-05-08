package mailman

import (
	"poste/util"
)

type Mailman struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var M = &Mailman{}

func (m *Mailman)Addr() string {
	return util.ToAddr(m.Host, m.Port)
}

func Serve(host string, port int) {
	handle(host, port)
}