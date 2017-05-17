package consul

import (
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func GetClient() *api.Client {

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		util.LogPanic("get consul client failed. error : %s", err)
	}
	return client
}

func GetKV() *api.KV {
	return GetClient().KV()
}

func GetHealth() *api.Health {
	return GetClient().Health()
}

func GetAgent() *api.Agent {
	return GetClient().Agent()
}
