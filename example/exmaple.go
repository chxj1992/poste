package main

import (
	"log"
	"net/http"
)

const (
	ID string = "000001"
	SECRET string = "example"
)

func main() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

	})
	log.Fatal(http.ListenAndServe(":0", nil))
}