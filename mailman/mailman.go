package mailman

import (
	"poste/util"
	"poste/consul"
	"github.com/serialx/hashring"
	"net"
	"math"
)

type Mailman struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

var (
	M = &Mailman{}
	mailmen []string
	mailmenRing *hashring.HashRing
)

func (m *Mailman)Addr() string {
	return util.ToAddr(m.Host, m.Port)
}

func OnShutDown() {
	util.LogInfo("mailman is shutting down ...")
	consul.Deregister(consul.Mailman, M.Host, M.Port)
	Retire()
	util.LogInfo("done!")
}

func Retire() {
	count := 0
	for _, hub := range userHubs {
		for _, c := range hub {
			c.disconnect(true)
		}
		count += 1
	}
	util.LogInfo("%d connections on the node are disconnected", count)
}

var mailmanCallback = func(values []*Mailman) {
	prevMailmen := mailmen
	mailmen = []string{}
	var newMailman string
	for _, m := range values {
		mailmen = append(mailmen, m.Addr())
		if len(prevMailmen) != 0 && !util.InSlice(m.Addr(), prevMailmen) {
			newMailman = m.Addr()
		}
	}
	mailmenRing = hashring.New(mailmen)

	if newMailman != "" && newMailman != M.Addr() {
		util.LogInfo("a new mailmen server registered : %s", newMailman)
		newMailmanPosition, _ := mailmenRing.GetNodePos(newMailman)
		currentPosition, _ := mailmenRing.GetNodePos(M.Addr())
		currentDistance := math.Abs(float64(newMailmanPosition - currentPosition))

		count := 0
		for _, m := range mailmen {
			position, _ := mailmenRing.GetNodePos(m)
			distance := math.Abs(float64(newMailmanPosition - position))
			if  distance != 0 && distance < currentDistance {
				count += 1
			}
			if count >= 2 {
				return
			}
		}

		util.LogInfo("current mailman is next to the new one")
		Retire()
	}
}

func beforeServe(addr *net.TCPAddr) {
	M.Host = addr.IP.String()
	M.Port = addr.Port

	go Watch(mailmanCallback)
}

func Serve(host string, port int) {
	defer func() {
		OnShutDown()
	}()

	handle(host, port)
}