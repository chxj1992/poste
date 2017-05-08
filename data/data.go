package data

import (
	"encoding/json"
	"poste/util"
)

type Data struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

func (d Data)Marshal() (bytes []byte) {
	bytes, err := json.Marshal(d)
	if err != nil {
		util.LogError("data marshal error : %s", err)
		return []byte{}
	}
	return
}
