package dispather

import (
	"log"
	"poste/mailman"
	"net"
	"net/http"
	"poste/util"
	"poste/consul"
	"github.com/serialx/hashring"
)

type Dispatcher struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var (
	D = &Dispatcher{}
	mailmenWs []string
	mailmenTcp []string
	mailmenWsRing *hashring.HashRing
	mailmenTcpRing *hashring.HashRing
)

func (d *Dispatcher)Addr() string {
	return util.ToAddr(d.Host, d.Port)
}

var callback = func(values []*mailman.Mailman) {
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
	log.Printf("[INFO] ws mailmen %s", mailmenWs)
	log.Printf("[INFO] tcp mailmen %s", mailmenTcp)

	mailmenWsRing = hashring.New(mailmenWs)
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
