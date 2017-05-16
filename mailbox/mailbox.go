package mailbox

import (
	"poste/data"
	"github.com/kr/beanstalk"
	"time"
	"poste/consul"
	"poste/util"
	"poste/ticket"
)

func Send(appId string, userId string, message string) {
	target := ticket.GetTicket(userId, appId, false)
	if target == "" {
		util.LogError("target is not connected")
		return
	}
	d := data.Data{Target:target, Message:message}
	bytes := d.Marshal()
	c := beanstalkClient()
	if c == nil {
		return
	}
	_, err := c.Put(bytes, 1, 0, time.Minute)
	if err != nil {
		util.LogError("beanstalk client send failed. error : %s", err)
	}
}

func beanstalkClient() *beanstalk.Conn {
	services := consul.Get(consul.Queue)
	if len(services) < 1 {
		util.LogError("queue service is not started")
		return nil
	}
	c, err := beanstalk.Dial("tcp", util.ToAddr(services[0].Host, services[0].Port))
	if err != nil {
		util.LogError("beanstalk client get failed. error : %s", err)
		return nil
	}

	return c
}
