package mailman

import (
	"github.com/hashicorp/consul/api"
	"poste/consul"
	"poste/util"
	"time"
)

var (
	mailmenValues []*Mailman
	update = make(chan int, 1)
	Refresh = make(chan int)
)

func Watch(callback func(mailmen []*Mailman)) {
	q := &api.QueryOptions{
		WaitTime: 60 * time.Minute,
	}
	update <- 1
	for {
		select {
		case <-update:
			util.LogInfo("update mailmen connection configuration from consul")
			updateFromConsul(q)
		case <-Refresh:
			util.LogInfo("refresh mailmen connection configuration")
		}
		callback(mailmenValues)
	}
}

func updateFromConsul(q *api.QueryOptions) {
	services, meta, err := consul.GetHealth().Service(string(consul.Mailman), "", false, q)
	q.WaitIndex = meta.LastIndex

	mailmenValues = []*Mailman{}
	if err != nil {
		util.LogError("consul get service failed. error : %s", err)
	}
	if services != nil {
		for _, s := range services {
			m := &Mailman{Host:s.Service.Address, Port:s.Service.Port}
			mailmenValues = append(mailmenValues, m)
		}
	}
	update <- 1
}
