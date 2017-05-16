package consul

import (
	"net"
	"net/http"
	"poste/util"
	"github.com/hashicorp/consul/api"
)

type ServiceType string

const (
	Dispatcher ServiceType = "dispatcher"
	Mailman ServiceType = "mailman"
	Api ServiceType = "api"
	Redis ServiceType = "redis"
	Queue ServiceType = "queue"
)

type Service struct {
	Name ServiceType
	Host string
	Port int
}

func RegisterServiceAndServe(serviceType ServiceType, host string, port int, tags []string, beforeServe func(addr *net.TCPAddr)) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	var err error
	address := util.ToAddr(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		util.LogError("%s%s server start listen failed: %s", serviceType, tags, err)
	}

	util.LogInfo("%s%s serves on %s", serviceType, tags, listener.Addr().String())
	addr := listener.Addr().(*net.TCPAddr)
	defer func() {
		Deregister(serviceType, addr.IP.String(), addr.Port)
	}()
	Register(serviceType, addr.IP.String(), addr.Port, tags, nil)

	if beforeServe != nil {
		beforeServe(addr)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		util.LogError("%s%s server start serve failed: %s", serviceType, tags, err)
	}
}

func Get(name ServiceType) []*Service {
	services, _, _ := GetHealth().Service(string(name), "", true, nil)

	values := []*Service{}
	for _, s := range services {
		values = append(values, &Service{
			Name: name,
			Host:s.Service.Address,
			Port:s.Service.Port,
		})
	}
	return values
}

func Clear(name ServiceType) {
	services, _, _ := GetHealth().Service(string(name), "", false, nil)
	for _, s := range services {
		GetAgent().ServiceDeregister(s.Service.ID)
	}
}

func Watch(name ServiceType, callback func(values []*api.ServiceEntry)) {
	q := &api.QueryOptions{}
	for {
		pairs, meta, err := GetHealth().Service(string(name), "", true, q)
		q.WaitIndex = meta.LastIndex

		values := []*api.ServiceEntry{}
		if err != nil {
			util.LogError("consul service get value failed. error : %s", err)
		}
		if pairs != nil {
			for _, pair := range pairs {
				values = append(values, pair)
			}
		}
		callback(values)
	}
}
