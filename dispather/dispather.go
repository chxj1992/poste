package dispather

import (
	"log"
	"time"
	"poste/mailman"
)

var mailmen []string

var callback = func(values []string) {
	mailmen = values
	log.Printf("[INFO] mailmen %s",mailmen)
}

func Serve() {
	go mailman.Watch(callback)
	for m := range mailmen {
		time.Sleep(time.Second * 5)
		log.Print(m)
	}
}
