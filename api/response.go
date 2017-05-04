package api

import (
	"encoding/json"
	"poste/util"
)

type Response struct {
	Err  string `json:"err"`
	Data string `json:"data"`
}

func (r Response)Marshal() (bytes []byte) {
	bytes, err := json.Marshal(r)
	if err != nil {
		util.LogError("response struct marshal error : %s", err)
		return []byte{}
	}
	return
}
