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

	mailmen []string
	mailmenRing *hashring.HashRing
)

func mailmanCallback(values []*mailman.Mailman) {
	mailmen = []string{}

	for _, m := range values {
		mailmen = append(mailmen, m.Addr())
	}

	util.LogInfo("mailmen %s", mailmen)
	mailmenRing = hashring.New(mailmen)
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