package mailman

import (
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"poste/util"
	"poste/consul"
	"poste/ticket"
	"poste/structure"
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
	uuid   string
}

var userHubs = map[string][]*Client{}

func handle(host string, port int) {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			util.LogError("serve websocket failed. error: %s", err)
			return
		}

		t := r.URL.Query().Get("ticket")
		info := ticket.GetUserInfo(t)

		if t != "" && (len(info) < 2 || info[0] == nil || info[1] == nil) {
			util.LogError("invalid ticket : %s", t)
			return
		}

		client := &Client{conn: conn, send: make(chan []byte, 256), ticket:t}

		if t != "" {
			client.uuid = ticket.UUID(info[0], info[1])
			userHubs[client.uuid] = append(userHubs[client.uuid], client)
			util.LogDebug("ticket: %s app: %s user: %s connected", t, info[1], info[0])
		}

		go readPump(client)
		writePump(client)
	})

	consul.RegisterServiceAndServe(consul.Mailman, host, port, nil, beforeServe)
}

func readPump(c *Client) {
	defer func() {
		c.disconnect()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				util.LogError("reading connection is unexpected closed : %s", err)
			}
			break
		}

		util.LogDebug("message received : %s", string(message))
		d := structure.Data{}
		err = json.Unmarshal(message, &d)
		if err != nil {
			util.LogError("invalid data structure : %s", string(message))
			continue
		}
		if len(userHubs[c.uuid]) == 0 {
			util.LogError("target %s not exists on this node", d.Target)
			continue
		}
		for _, t := range userHubs[c.uuid] {
			t.send <- []byte(d.Message)
		}
	}
}

func writePump(c *Client) {
	defer func() {
		c.disconnect()
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

func (c *Client) disconnect() {
	m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server")
	c.conn.WriteMessage(websocket.CloseMessage, m)

	if c.uuid == "" || len(userHubs[c.uuid]) == 0 {
		return
	}

	clients := []*Client{}
	for _, client := range userHubs[c.uuid] {
		if client != c {
			clients = append(clients)
		}
	}

	userHubs[c.uuid] = clients

	if len(clients) == 0 {
		util.LogDebug("user %s with ticket %s is offline", c.uuid, c.ticket)
	} else {
		util.LogDebug("user %s with ticket %s closed 1 connection", c.uuid, c.ticket)
	}
}