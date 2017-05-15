package mailbox

import (
	"testing"
)

func TestSend(t *testing.T) {
	Send("000001", "1", "hello world to user 1")
	//Send("000001", "2", "hello world to user 2")
}