package data

import "poste/mailman"

type Data struct {
	Target     string `json:"target"`
	ServerType mailman.ServerType `json:"type,omitempty"`
	Message    string `json:"message"`
}
