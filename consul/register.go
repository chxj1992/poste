package consul

import (
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func Register(name ServiceType, host string, port int, tags []string, check *api.AgentServiceCheck) error {
	util.LogInfo("%s service %s registered", name, util.ToAddr(host, port))
	s := &api.AgentServiceRegistration{
		ID:ServiceId(name, host, port),
		Name:string(name),
		Address:host,
		Port:port,
		Tags:tags,
	}
	if check == nil {
		s.Check = &api.AgentServiceCheck{
			HTTP: "http://" + util.ToAddr(host, port),
			Interval: "10s",
		}
	} else {
		s.Check = check
	}
	return GetAgent().ServiceRegister(s)
}

func Deregister(name ServiceType, host string, port int) error {
	util.LogInfo("%s service %s is deregistered for consul", name, util.ToAddr(host, port))
	return GetAgent().ServiceDeregister(ServiceId(name, host, port))
}

func ServiceId(name ServiceType, host string, port int) string {
	addr := util.ToAddr(host, port)
	sum := util.Md5(addr)
	return string(name) + "_" + sum
}