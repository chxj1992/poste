package dispather

import (
	"poste/consul"
	"log"
	"github.com/hashicorp/consul/api"
	"poste/data"
	"github.com/kr/beanstalk"
	"time"
	"encoding/json"
	"poste/mailman"
	"github.com/gorilla/websocket"
)

func Consume() {
	var count = 0
	consul.KVWatch("queue", func(values []*api.KVPair) {
		queues := getQueues(values)

		// only keep the latest consumers loop
		var work = make(chan int, 1)
		count += 1
		work <- count
		log.Printf("watching queue. %s", values)
		go func() {
			for {
				v := <-work
				if v != count {
					log.Printf("v %s c %s", v, count)
					break;
				} else {
					work <- v
					consuming(queues)
				}
			}
		}()
	})
}

func getQueues(values []*api.KVPair) []*beanstalk.Conn {
	queues := []*beanstalk.Conn{}
	for _, pair := range values {
		c, err := beanstalk.Dial("tcp", string(pair.Value))
		if err != nil {
			log.Printf("get queue failed. error %s", err)
			consul.KVDelete(string(pair.Key))
			continue
		}
		queues = append(queues, c)
	}
	return queues
}

func consuming(queues []*beanstalk.Conn) {
	var d data.Data
	for _, c := range queues {
		id, body, err := c.Reserve(5 * time.Second)
		if err != nil {
			continue
		}
		c.Delete(id)
		log.Printf("get queue. ID : %s . data : %s", id, body)
		json.Unmarshal(body, &d)

		if d.ServerType == mailman.WsType {
			addr, ok := mailmenWsRing.GetNode(d.Target)
			if !ok {
				log.Printf("[ERROR] get wsmailman failed. addr : %s", addr)
			}
			mailmenWsClients[addr].WriteMessage(websocket.TextMessage, d.Marshal())
		}
	}
}
