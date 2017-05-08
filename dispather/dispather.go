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

	mailmenWs []string
	mailmenWsClients map[string]*websocket.Conn
	mailmenWsRing *hashring.HashRing

)

var callback = func(values []*mailman.Mailman) {
	mailmenWs = []string{}
	mailmenWsClients = map[string]*websocket.Conn{}

	util.LogInfo("values %s", values)
	// establish connection with every mailman server
	for _, m := range values {

		if m.ServerType == mailman.WsType {
			mailmenWs = append(mailmenWs, m.Addr())
			c, _, err := websocket.DefaultDialer.Dial("ws://" + m.Addr() + "/connect", nil)
			if err != nil {
				util.LogError("connect to mailman failed : %s", err)
				continue
			}
			util.LogInfo("connected to ws mailman : %s", m.Addr())
			mailmenWsClients[m.Addr()] = c
		}
	}

	util.LogInfo("ws mailmen %s", mailmenWs)
	mailmenWsRing = hashring.New(mailmenWs)
}

func Serve(host string, port int) {
	go mailman.Watch(callback)
	go Consume()

	consul.RegisterServiceAndServe("dispatcher", host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	D.Host = addr.IP.String()
	D.Port = addr.Port
}