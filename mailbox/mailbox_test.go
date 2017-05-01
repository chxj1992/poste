package mailbox

import (
	"testing"
	"poste/mailman"
)

func TestSend(t *testing.T) {
	Send("1", "hello world", mailman.WsType)
}