package ticket

import (
	"testing"
	"log"
)


func TestBind(t *testing.T) {
	log.Print(GetTicket("111", "0001", true))
}

func TestGetUserIdByTicket(t *testing.T) {
	log.Print(GetUserInfo("849d3290efe17bcb8a022d2933926b76"))
}