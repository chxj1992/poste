package config

import (
	"poste/consul"
	"runtime"
	"path"
	"github.com/jinzhu/configor"
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func Init() {
	registerRedis()
	registerQueue()
}

func registerRedis() {
	consul.Clear(consul.Redis)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/redis.json")
	redisConf := consul.Service{}
	configor.Load(&redisConf, configPath)

	consul.Register(consul.Redis, redisConf.Host, redisConf.Port, nil, &api.AgentServiceCheck{
		TCP: util.ToAddr(redisConf.Host, redisConf.Port),
		Interval: "10s",
	})
}

func registerQueue() {
	consul.Clear(consul.Queue)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/queues.json")
	queueConfigs := []consul.Service{}
	configor.Load(&queueConfigs, configPath)

	for _, queueConfig := range queueConfigs {
		consul.Register(consul.Queue, queueConfig.Host, queueConfig.Port, nil, &api.AgentServiceCheck{
			TCP: util.ToAddr(queueConfig.Host, queueConfig.Port),
			Interval: "10s",
		})
	}
}