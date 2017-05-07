package mailbox

import (
	"testing"
	"poste/mailman"
)

func TestSend(t *testing.T) {
	Send("000001", "1", "hello world to user 1", mailman.WsType)
	Send("000001", "2", "hello world to user 2", mailman.WsType)
}