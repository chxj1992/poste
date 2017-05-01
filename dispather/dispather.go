package dispather

import (
	"log"
	"poste/mailman"
	"net"
	"net/http"
	"poste/util"
	"poste/consul"
	"github.com/serialx/hashring"
	"github.com/gorilla/websocket"
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

	mailmenTcp []string
	mailmenTcpRing *hashring.HashRing
)

func (d *Dispatcher)Addr() string {
	return util.ToAddr(d.Host, d.Port)
}

var callback = func(values []*mailman.Mailman) {
	mailmenWs = []string{}
	mailmenWsClients = map[string]*websocket.Conn{}
	mailmenTcp = []string{}

	for _, m := range values {
		if m.ServerType == mailman.WsType {
			mailmenWs = append(mailmenWs, m.Addr())
			c, _, err := websocket.DefaultDialer.Dial("ws://" + m.Addr(), nil)
			if err != nil {
				log.Printf("connect to mailman failed : %s", err)
			}
			mailmenWsClients[m.Addr()] = c
		}
		if m.ServerType == mailman.TcpType {
			mailmenTcp = append(mailmenTcp, m.Addr())
		}
	}

	log.Printf("[INFO] ws mailmen %s", mailmenWs)
	mailmenWsRing = hashring.New(mailmenWs)

	log.Printf("[INFO] tcp mailmen %s", mailmenTcp)
	mailmenTcpRing = hashring.New(mailmenTcp)
}

func Serve(host string, port int) {
	go mailman.Watch(callback)
	go Consume()

	handleHttp()

	address := util.ToAddr(host, port)

	var err error
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("dispatcher server start failed: ", err)
	}
	log.Printf("dispatcher serves on %s", listener.Addr().String())
	addr := listener.Addr().(*net.TCPAddr)
	defer func() {
		consul.Deregister("dispatcher", addr.IP.String(), addr.Port)
	}()
	consul.Register("dispatcher", addr.IP.String(), addr.Port, nil)
	D.Host = addr.IP.String()
	D.Port = addr.Port

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("dispatcher server start failed: ", err)
	}
}
func handleHttp() {
	http.HandleFunc("/mailmen/ws", MailmenWs)
	http.HandleFunc("/mailmen/tcp", MailmenTcp)
}
