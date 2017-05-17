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

	util.LogDebug("mailmen : %s", mailmen)
	mailmenRing = hashring.New(mailmen)
}

func OnShutDown() {
	util.LogInfo("api service is shutting down ...")
	consul.Deregister(consul.Api, A.Host, A.Port)
	util.LogInfo("retired!")
}

func Serve(host string, port int) {
	defer func() {
		OnShutDown()
	}()

	go mailman.Watch(mailmanCallback)

	oauthSvr = buildSrv()
	handleRequest()

	consul.RegisterServiceAndServe("api", host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	A.Host = addr.IP.String()
	A.Port = addr.Port
}