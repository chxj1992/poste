package mailman

import (
	"testing"
	"log"
)

func TestWatch(t *testing.T) {
	callback := func(values []string) {
		log.Print(values)
	}
	Watch(callback)
}

func TestServe(t *testing.T)  {
	Serve("127.0.0.1", 0, WsType)
}