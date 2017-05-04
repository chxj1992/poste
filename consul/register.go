package consul

import (
	"github.com/hashicorp/consul/api"
	"crypto/md5"
	"poste/util"
	"io"
	"fmt"
)

func Register(name ServiceType, host string, port int, tags []string) error {
	util.LogInfo("%s service %s registered", name, util.ToAddr(host, port))
	s := &api.AgentServiceRegistration{
		ID:ServiceId(name, host, port),
		Name:string(name),
		Address:host,
		Port:port,
		Tags:tags,
	}
	return GetAgent().ServiceRegister(s)
}

func Deregister(name ServiceType, host string, port int) error {
	util.LogInfo("%s service %s deregistered", name, util.ToAddr(host, port))
	return GetAgent().ServiceDeregister(ServiceId(name, host, port))
}

func ServiceId(name ServiceType, host string, port int) string {
	addr := util.ToAddr(host, port)
	h := md5.New()
	io.WriteString(h, addr)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return string(name) + "_" + sum
}