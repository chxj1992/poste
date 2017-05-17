package register

import (
	"poste/consul"
	"runtime"
	"path"
	"github.com/jinzhu/configor"
	"github.com/hashicorp/consul/api"
	"poste/util"
)

func Init() {
	initRedis()
	initQueue()
}

func initRedis() {
	consul.Clear(consul.Redis)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/redis.json")
	redisConf := consul.Service{}
	configor.Load(&redisConf, configPath)

	registerRedis(redisConf.Host, redisConf.Port)
}

func initQueue() {
	consul.Clear(consul.Queue)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/queues.json")
	queueConfigs := []consul.Service{}
	configor.Load(&queueConfigs, configPath)

	for _, queueConfig := range queueConfigs {
		registerQueue(queueConfig.Host, queueConfig.Port)
	}
}

func registerRedis(host string, port int) {
	consul.Register(consul.Redis, host, port, nil, &api.AgentServiceCheck{
		TCP: util.ToAddr(host, port),
		Interval: "10s",
	})
}

func registerQueue(host string, port int) {
	consul.Register(consul.Queue, host, port, nil, &api.AgentServiceCheck{
		TCP: util.ToAddr(host, port),
		Interval: "10s",
	})
}