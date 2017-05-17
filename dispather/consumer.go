package dispather

import (
	"poste/consul"
	"github.com/hashicorp/consul/api"
	"poste/data"
	"github.com/kr/beanstalk"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"poste/util"
	"strings"
	"poste/ticket"
	"poste/mailman"
)

func Consume() {
	var count = 0
	consul.Watch("queue", func(values []*api.ServiceEntry) {
		queues := getQueues(values)

		// only keep the latest consumers loop
		var work = make(chan int, 1)
		count += 1
		work <- count
		util.LogDebug("watching queue. %s", values)
		go func() {
			for {
				v := <-work
				if v != count {
					util.LogDebug("consumer version: %s, current version: %s", v, count)
					break;
				} else {
					work <- v
					consuming(queues)
				}
			}
		}()
	})
}

func getQueues(values []*api.ServiceEntry) []*beanstalk.Conn {
	queues := []*beanstalk.Conn{}
	for _, service := range values {
		c, err := beanstalk.Dial("tcp", util.ToAddr(service.Service.Address, service.Service.Port))
		if err != nil {
			util.LogError("get queue failed. error %s", err)
			consul.Deregister(consul.Queue, service.Service.Address, service.Service.Port)
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
		util.LogDebug("get data from queue. ID : %s . data : %s", id, body)
		json.Unmarshal(body, &d)

		t := util.Base64Decode(d.Target)
		info := strings.Split(t, ticket.SepChar)

		if mailmenRing == nil {
			time.Sleep(time.Second)
			continue
		}

		addr, ok := mailmenRing.GetNode(info[0])
		if !ok {
			util.LogError("get mailman failed. addr : %s", addr)
		}
		if mailmenClients[addr] != nil {
			err := mailmenClients[addr].WriteMessage(websocket.TextMessage, d.Marshal())
			if err != nil {
				util.LogError("dispatch message failed: %s", err)
				mailman.Refresh <- 1
			} else {
				util.LogDebug("message sent : %s", string(d.Marshal()))
				c.Delete(id)
			}
		} else {
			util.LogError("addr %s mailman connection not exists", addr)
			mailman.Refresh <- 1
		}
	}
}
