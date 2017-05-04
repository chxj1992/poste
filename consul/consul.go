package consul

import (
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func GetClient() *api.Client {

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		util.LogError("consul kv get value failed. error : %s", err)
		panic("get consul client failed")
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
