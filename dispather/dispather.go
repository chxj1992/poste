package dispather

import (
	"log"
	"poste/mailman"
	"net"
	"net/http"
	"poste/util"
	"poste/consul"
)

type Dispatcher struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var (
	D = &Dispatcher{}
	mailmen []*mailman.Mailman
)

func (d *Dispatcher)Addr() string {
	return util.ToAddr(d.Host, d.Port)
}

var callback = func(values []*mailman.Mailman) {
	mailmen = values

	log.Printf("[INFO] mailmen %s", mailmen)
}

func Serve(host string, port int) {
	go mailman.Watch(callback)

	http.HandleFunc("/mailmen", Mailmen);

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
