package dispather

import (
	"poste/mailman"
	"github.com/serialx/hashring"
	"github.com/gorilla/websocket"
	"net"
	"poste/util"
	"poste/consul"
)

type Dispatcher struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var (
	D = &Dispatcher{}

	mailmen []string
	mailmenClients map[string]*websocket.Conn
	mailmenRing *hashring.HashRing
)

func OnShutDown() {
	util.LogInfo("dispather is shutting down ...")
	closeMailmanConn()
	consul.Deregister(consul.Dispatcher, D.Host, D.Port)
	util.LogInfo("done!")
}

func closeMailmanConn() {
	for _, c := range mailmenClients {
		m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close connection normally")
		c.WriteMessage(websocket.CloseMessage, m)
	}
}

var mailmanCallback = func(values []*mailman.Mailman) {
	closeMailmanConn()

	mailmen = []string{}
	mailmenClients = map[string]*websocket.Conn{}

	util.LogInfo("watching mailmen returned %s", values)

	// establish connection with every mailman server
	for _, m := range values {
		mailmen = append(mailmen, m.Addr())
		c, _, err := websocket.DefaultDialer.Dial("ws://" + m.Addr() + "/connect", nil)
		if err != nil {
			util.LogError("connect to mailman failed : %s", err)
			mailman.Refresh <- 1
			continue
		}
		util.LogInfo("connected to mailman : %s", m.Addr())
		mailmenClients[m.Addr()] = c
	}

	util.LogInfo("mailmen %s", mailmen)
	mailmenRing = hashring.New(mailmen)
}

func Serve(host string, port int) {
	defer func() {
		OnShutDown()
	}()

	go mailman.Watch(mailmanCallback)
	go Consume()

	consul.RegisterServiceAndServe("dispatcher", host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	D.Host = addr.IP.String()
	D.Port = addr.Port
}