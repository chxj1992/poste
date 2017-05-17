package mailman

import (
	"github.com/hashicorp/consul/api"
	"poste/consul"
	"poste/util"
)

var (
	mailmenValues []*Mailman
	update = make(chan int, 1)
	Refresh = make(chan int)
	startup = true
)

func Watch(callback func(mailmen []*Mailman)) {
	q := &api.QueryOptions{}
	update <- 1
	for {
		select {
		case <-update:
			util.LogDebug("update mailmen connection configuration from consul")
			updateFromConsul(q)
		case <-Refresh:
			util.LogDebug("refresh mailmen connection configuration")
		}
		go callback(mailmenValues)
		startup = false
	}
}

func updateFromConsul(q *api.QueryOptions) {
	services, meta, err := consul.GetHealth().Service(string(consul.Mailman), "", true, q)
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

	if !startup {
		update <- 1
	}
}
