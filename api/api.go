package api

import (
	"gopkg.in/oauth2.v3/server"
	"net"
	"poste/consul"
	"poste/mailman"
	"poste/util"
	"github.com/serialx/hashring"
)

var oauthSvr *server.Server

type Api struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var
(
	A = &Api{}

	mailmenWs []string
	mailmenWsRing *hashring.HashRing

	//TODO:tcp mailman server
	mailmenTcp []string
	mailmenTcpRing *hashring.HashRing
)
func mailmanCallback(values []*mailman.Mailman) {
	mailmenWs = []string{}
	mailmenTcp = []string{}

	for _, m := range values {
		if m.ServerType == mailman.WsType {
			mailmenWs = append(mailmenWs, m.Addr())
		}
		if m.ServerType == mailman.TcpType {
			mailmenTcp = append(mailmenTcp, m.Addr())
		}
	}

	util.LogInfo("ws mailmen %s", mailmenWs)
	mailmenWsRing = hashring.New(mailmenWs)

	util.LogInfo("tcp mailmen %s", mailmenTcp)
	mailmenTcpRing = hashring.New(mailmenTcp)
}

func Serve(host string, port int) {
	oauthSvr = buildSrv()

	go mailman.Watch(mailmanCallback)

	handleRequest()

	consul.RegisterServiceAndServe("api", host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	A.Host = addr.IP.String()
	A.Port = addr.Port
}