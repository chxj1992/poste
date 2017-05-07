package mailman

import (
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"poste/data"
	"poste/util"
	"poste/consul"
	"poste/ticket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var hub = map[string][]*Client{}

func serveWs(host string, port int) {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			util.LogError("serve websocket failed. error: %s", err)
			return
		}
		client := &Client{conn: conn, send: make(chan []byte, 256)}

		if t := r.URL.Query().Get("ticket"); t != "" {
			info := ticket.GetUserInfo(t)
			if info[0] == nil {
				util.LogError("ticket invalid : %s", t)
				return
			}
			hub[t] = append(hub[t], client)
			util.LogInfo("ticket: %s app: %s user: %s connected", t, info[1], info[0])
		}
		go readWs(client)
		writeWs(client)
	})

	consul.RegisterServiceAndServe(consul.Mailman, host, port, []string{string(WsType)}, beforeServe)
}

func readWs(c *Client) {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				util.LogError("read ws error: %s", err)
			}
			break
		}
		util.LogInfo("message received : %s", string(message))
		d := data.Data{}
		json.Unmarshal(message, &d)
		for _, t := range hub[d.Target] {
			t.send <- []byte(d.Message)
		}
	}
}

func writeWs(c *Client) {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			util.LogError("get websocket writer failed. error : %s", err)
			return
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			util.LogError("write message to ws connection failed. error : %s", err)
			return
		}
	}
}