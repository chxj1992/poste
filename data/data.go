package data

import (
	"encoding/json"
	"log"
)

type ServerType string

type Data struct {
	Target     string `json:"target"`
	ServerType ServerType `json:"type"`
	Message    string `json:"message"`
}

func (d Data)Marshal() []byte {
	bytes,err := json.Marshal(d)
	if err != nil {
		log.Printf("data marshal error : %s", err)
		return []byte{}
	}
	return bytes
}
