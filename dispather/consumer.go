package dispather

import (
	"poste/consul"
	"log"
	"github.com/hashicorp/consul/api"
	"poste/data"
	"github.com/kr/beanstalk"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"poste/util"
	"strings"
	"poste/ticket"
)

func Consume() {
	var count = 0
	consul.KVWatch("queue", func(values []*api.KVPair) {
		queues := getQueues(values)

		// only keep the latest consumers loop
		var work = make(chan int, 1)
		count += 1
		work <- count
		util.LogInfo("watching queue. %s", values)
		go func() {
			for {
				v := <-work
				if v != count {
					log.Printf("consumer version: %s, current version: %s", v, count)
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
			util.LogError("get queue failed. error %s", err)
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
		util.LogInfo("get data from queue. ID : %s . data : %s", id, body)
		json.Unmarshal(body, &d)

		t := util.Base64Decode(d.Target)
		info := strings.Split(t, ticket.SepChar)
		addr, ok := mailmenRing.GetNode(info[0])
		if !ok {
			util.LogError("get mailman failed. addr : %s", addr)
		}
		if mailmenClients[addr] != nil {
			err := mailmenClients[addr].WriteMessage(websocket.TextMessage, d.Marshal())
			if err != nil {
				util.LogError("dispatch message failed: %s", err)
			}
		} else {
			util.LogError("addr %s mailman connection hub not exists", addr)
		}
	}
}
