package mailman

import (
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"poste/data"
	"poste/util"
	"poste/consul"
	"poste/ticket"
	"net"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	ticket string
}

var hub = map[string][]*Client{}

func handle(host string, port int) {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			util.LogError("serve websocket failed. error: %s", err)
			return
		}

		t := r.URL.Query().Get("ticket")
		info := ticket.GetUserInfo(t)

		util.LogInfo("%s %s", t, info)

		if t != "" && (len(info) < 2 || info[0] == nil || info[1] == nil) {
			util.LogError("invalid ticket : %s", t)
			return
		}

		client := &Client{conn: conn, send: make(chan []byte, 256), ticket:t}

		if t != "" {
			hub[t] = append(hub[t], client)
			util.LogInfo("ticket: %s app: %s user: %s connected", t, info[1], info[0])
		}

		go readPump(client)
		writePump(client)
	})

	consul.RegisterServiceAndServe(consul.Mailman, host, port, nil, beforeServe)
}

func beforeServe(addr *net.TCPAddr) {
	M.Host = addr.IP.String()
	M.Port = addr.Port
}

func readPump(c *Client) {
	defer func() {
		removeFromHub(c)
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				util.LogError("reading connection is unexpected closed : %s", err)
			}
			break
		}

		util.LogInfo("message received : %s", string(message))
		d := data.Data{}
		json.Unmarshal(message, &d)

		if len(hub[d.Target]) == 0 {
			util.LogError("target %s not exists on this node", d.Target)
			break
		}
		for _, t := range hub[d.Target] {
			t.send <- []byte(d.Message)
		}
	}
}

func writePump(c *Client) {
	defer func() {
		removeFromHub(c)
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			util.LogError("get message from sending channel failed.")
			break
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			util.LogError("get websocket writer failed. error : %s", err)
			break
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			util.LogError("write message to ws connection failed. error : %s", err)
			break
		}
	}
}

func removeFromHub(c *Client) {

	m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server")
	c.conn.WriteMessage(websocket.CloseMessage, m)

	if c.ticket == "" || len(hub[c.ticket]) == 0 {
		return
	}

	clients := []*Client{}
	for _, client := range hub[c.ticket] {
		if client != c {
			clients = append(clients)
		}
	}

	hub[c.ticket] = clients

	if len(clients) == 0 {
		util.LogInfo("target %s is offline", c.ticket)
	} else {
		util.LogInfo("target %s closed 1 connection", c.ticket)
	}
}