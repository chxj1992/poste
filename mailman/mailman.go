package mailman

import (
	"poste/util"
	"poste/consul"
)

type Mailman struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var M = &Mailman{}

func (m *Mailman)Addr() string {
	return util.ToAddr(m.Host, m.Port)
}

func OnShutDown() {
	util.LogInfo("mailman is shutting down ...")
	consul.Deregister(consul.Mailman, M.Host, M.Port)
	util.LogInfo("done!")
}

func Serve(host string, port int) {
	defer func() {
		OnShutDown()
	}()

	handle(host, port)
}