package ticket

import (
	"github.com/go-redis/redis"
	"poste/util"
	"time"
	"poste/consul"
)

const SepChar = "#"

func Client() *redis.Client {
	services := consul.Get(consul.Redis)
	if len(services) <= 0 {
		util.LogPanic("redis is not currently available, try `poste init` to initialize the services from config.")
	}
	service := services[0]
	redisConf := redis.Options{
		Addr: util.ToAddr(service.Host, service.Port),
	}

	return redis.NewClient(&redisConf)
}

func UUID(userId string, appId string) string {
	return util.Md5(userId + "-" + appId)
}

func GetTicket(userId string, appId string, refresh bool) (ticket string) {
	client := Client()
	if refresh == false {
		return client.Get("ticket:" + appId + ":" + userId).Val()
	}

	ticket = util.Base64Encode(UUID(userId, appId) + SepChar + util.RandStr(8))
	client.Set("ticket:" + appId + ":" + userId, ticket, 0)

	client.Pipelined(func(pipe redis.Pipeliner) error {
		client.HMSet("user:" + ticket, map[string]interface{}{"userId":userId, "appId":appId})
		client.Expire("user:" + ticket, 24 * time.Hour)
		return nil
	})
	return
}

func GetUserInfo(ticket string) []interface{} {
	client := Client()
	return client.HMGet("user:" + ticket, "userId", "appId").Val()
}
