package data

import "poste/mailman"

type Data struct {
	Target     string `json:"target"`
	ServerType mailman.ServerType `json:"type"`
	Message    string `json:"message"`
}
