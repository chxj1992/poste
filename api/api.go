package api

import (
	"gopkg.in/oauth2.v3/server"
	"net"
	"poste/consul"
)

var oauthSvr *server.Server

type Api struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var A = &Api{}

func Serve(host string, port int) {
	oauthSvr = buildSrv()

	handleRequest()

	consul.RegisterServiceAndServe("api", host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	A.Host = addr.IP.String()
	A.Port = addr.Port
}