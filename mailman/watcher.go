package mailman

import (
	"github.com/hashicorp/consul/api"
	"log"
	"poste/consul"
	"poste/util"
)

func Watch(callback func(values []string)) {
	q := &api.QueryOptions{}
	for ; ; {
		services, meta, err := consul.GetHealth().Service("mailman", "", false, q)
		q.WaitIndex = meta.LastIndex

		values := []string{}
		if err != nil {
			log.Printf("[ERROR] consul kv get value failed. error : %s", err)
		}
		if services != nil {
			for _, service := range services {
				addr := util.ToAddr(service.Service.Address, service.Service.Port)
				values = append(values, addr)
			}
		}
		callback(values)
	}
}
