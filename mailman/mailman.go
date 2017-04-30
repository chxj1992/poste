package mailman

type ServerType string

const (
	TcpType ServerType = "tcp"
	WsType ServerType = "ws"
)

var  (
	Host string
	Port int
)

func Serve(host string, port int, serverType ServerType) {
	if serverType == WsType {
		serveWs(host, port)
	}
}

