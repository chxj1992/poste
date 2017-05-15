package consul

import (
	"net"
	"net/http"
	"poste/util"
)

type ServiceType string

const (
	Dispatcher ServiceType = "dispatcher"
	Mailman ServiceType = "mailman"
	Api ServiceType = "api"
	Redis ServiceType = "redis"
	Queue ServiceType = "queue"
)

func RegisterServiceAndServe(serviceType ServiceType, host string, port int, tags []string, beforeServe func(addr *net.TCPAddr)) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	var err error
	address := util.ToAddr(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		util.LogInfo("%s%s server start listen failed: %s", serviceType, tags, err)
	}

	util.LogInfo("%s%s serves on %s", serviceType, tags, listener.Addr().String())
	addr := listener.Addr().(*net.TCPAddr)
	defer func() {
		Deregister(serviceType, addr.IP.String(), addr.Port)
	}()
	Register(serviceType, addr.IP.String(), addr.Port, tags)

	if beforeServe != nil {
		beforeServe(addr)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		util.LogInfo("%s%s server start serve failed: %s", serviceType, tags, err)
	}
}