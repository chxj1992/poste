package dispather

import (
	"poste/consul"
	"log"
	"github.com/hashicorp/consul/api"
	"poste/data"
	"github.com/kr/beanstalk"
	"time"
	"encoding/json"
)

var message data.Data
var count = 0

func Consume() {
	consul.KVWatch("queue", func(values []*api.KVPair) {
		var queue = make(chan int, 1)
		count += 1
		queue <- count
		log.Printf("watching queue. %s", values)
		go func() {
			for {
				v := <-queue
				if v != count {
					log.Printf("v %s c %s", v, count)
					break;
				} else {
					queue <- v
					consuming(values)
				}
			}
		}()
	})
}

func consuming(values []*api.KVPair) {
	for _, pair := range values {
		c, err := beanstalk.Dial("tcp", string(pair.Value))
		if err != nil {
			log.Printf("get queue failed. error %s", err)
			consul.KVDelete(string(pair.Key))
			continue
		}
		id, body, err := c.Reserve(5 * time.Second)
		if err != nil {
			log.Printf("get message failed. error %s", err)
			continue
		}
		log.Printf("get queue. ID : %s . data : %s", id, body)
		json.Unmarshal(body, &message)
		log.Printf("message : %s", message)
	}
}
