package consul

import (
	"github.com/hashicorp/consul/api"
	"crypto/md5"
	"log"
	"poste/util"
	"io"
	"fmt"
)

func Register(name string, host string, port int, tags []string) error {
	log.Printf("%s service %s registered", name, util.ToAddr(host, port))
	service := &api.AgentServiceRegistration{
		ID:ServiceId(name, host, port),
		Name:name,
		Address:host,
		Port:port,
		Tags:tags,
	}
	return GetAgent().ServiceRegister(service)
}

func Deregister(name string, host string, port int) error {
	log.Printf("%s service %s deregistered", name, util.ToAddr(host, port))
	return GetAgent().ServiceDeregister(ServiceId(name, host, port))
}

func ServiceId(name string, host string, port int) string {
	addr := util.ToAddr(host, port)
	h := md5.New()
	io.WriteString(h, addr)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return name + "_" + sum
}