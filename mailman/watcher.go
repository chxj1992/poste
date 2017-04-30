package mailman

import (
	"github.com/hashicorp/consul/api"
	"log"
	"poste/consul"
)

func Watch(callback func(mailmen []*Mailman)) {
	q := &api.QueryOptions{}
	for ; ; {
		services, meta, err := consul.GetHealth().Service("mailman", "", false, q)
		q.WaitIndex = meta.LastIndex

		values := []*Mailman{}
		if err != nil {
			log.Printf("[ERROR] consul kv get value failed. error : %s", err)
		}
		if services != nil {
			for _, service := range services {
				m := &Mailman{Host:service.Service.Address, Port:service.Service.Port, ServerType:ServerType(service.Service.Tags[0])}
				values = append(values, m)
			}
		}
		callback(values)
	}
}
