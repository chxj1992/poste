package mailman

import (
	"github.com/hashicorp/consul/api"
	"poste/consul"
	"poste/util"
)

func Watch(callback func(mailmen []*Mailman)) {
	q := &api.QueryOptions{}
	for {
		services, meta, err := consul.GetHealth().Service(string(consul.Mailman), "", false, q)
		q.WaitIndex = meta.LastIndex

		values := []*Mailman{}
		if err != nil {
			util.LogError("consul get service failed. error : %s", err)
		}
		if services != nil {
			for _, s := range services {
				m := &Mailman{Host:s.Service.Address, Port:s.Service.Port}
				values = append(values, m)
			}
		}
		callback(values)
	}
}
