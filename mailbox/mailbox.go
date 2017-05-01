package mailbox

import (
	"poste/data"
	"github.com/kr/beanstalk"
	"log"
	"time"
	"poste/consul"
	"poste/util"
)

func Send(target string, message string, serverType data.ServerType) {
	d := data.Data{Target:target, Message:message, ServerType:serverType}
	bytes := d.Marshal()
	c := beanstalkClient()
	_, err := c.Put(bytes, 1, 0, time.Minute)
	if err != nil {
		log.Printf("beanstalk client send failed. error : %s", err)
	}
}

func beanstalkClient() *beanstalk.Conn {
	values := consul.KVValues("queue")
	c, err := beanstalk.Dial("tcp", util.Random(values))
	if err != nil {
		log.Printf("beanstalk client get failed. error : %s", err)
		panic(err)
	}
	return c
}
