package ticket

import (
	"github.com/go-redis/redis"
	"poste/util"
	"github.com/jinzhu/configor"
	"path"
	"runtime"
	"time"
)

const SepChar = "#"

func Client() *redis.Client {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/redis.json")
	redisConf := redis.Options{}
	configor.Load(&redisConf, configPath)
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
