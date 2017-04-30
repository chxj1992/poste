package mailman

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net"
	"poste/consul"
	"poste/util"
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

func serveWs(host string, port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("serve websocket failed. error: %s", err)
			return
		}
		client := &Client{conn: conn, send: make(chan []byte, 256)}
		go readWs(client)
		writeWs(client)
	})
	address := util.ToAddr(host, port)

	var err error
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("mailman server start failed: ", err)
	}
	log.Printf("websocket mailman serves on %s", listener.Addr().String())
	addr := listener.Addr().(*net.TCPAddr)
	defer func() {
		consul.Deregister("mailman", addr.IP.String(), addr.Port)
	}()
	consul.Register("mailman", addr.IP.String(), addr.Port, []string(string(WsType)))
	Host = addr.IP.String()
	Port = addr.Port

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("mailman server start failed: ", err)
	}
}

func readWs(c *Client) {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("read ws error: %s", err)
			}
			break
		}
		log.Printf("message received : %s", string(message))
		c.send <- []byte("message received : " + string(message))
	}
}

func writeWs(c *Client) {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}