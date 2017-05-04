package data

import (
	"encoding/json"
	"poste/util"
)

type ServerType string

type Data struct {
	Target     string `json:"target"`
	ServerType ServerType `json:"type"`
	Message    string `json:"message"`
}

func (d Data)Marshal() (bytes []byte) {
	bytes, err := json.Marshal(d)
	if err != nil {
		util.LogError("data marshal error : %s", err)
		return []byte{}
	}
	return
}
