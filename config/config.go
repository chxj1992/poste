package config

import (
	"poste/consul"
	"runtime"
	"path"
	"github.com/jinzhu/configor"
	"net"
	"github.com/hashicorp/consul/api"
	"strconv"
)

func Init() {
	registerRedis()
	registerQueue()
}

func registerRedis() {
	consul.Clear(consul.Redis)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/redis.json")
	redisConf := map[string]string{}
	configor.Load(&redisConf, configPath)

	host, port, _ := net.SplitHostPort(redisConf["Addr"])
	p, _ := strconv.Atoi(port)

	consul.Register(consul.Redis, host, p, nil, &api.AgentServiceCheck{
		TCP: redisConf["Addr"],
		Interval: "10s",
	})
}

func registerQueue() {
	consul.Clear(consul.Queue)

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/queues.json")
	queueConfigs := []map[string]string{}
	configor.Load(&queueConfigs, configPath)

	for _, queueConfig := range queueConfigs {
		host, port, _ := net.SplitHostPort(queueConfig["Addr"])
		p, _ := strconv.Atoi(port)
		consul.Register(consul.Queue, host, p, nil, &api.AgentServiceCheck{
			TCP: queueConfig["Addr"],
			Interval: "10s",
		})
	}
}