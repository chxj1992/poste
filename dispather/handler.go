package dispather

import (
	"encoding/json"
	"net/http"
)

func MailmenWs(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(mailmenWs)
	w.Write(bytes)
}

func MailmenTcp(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(mailmenTcp)
	w.Write(bytes)
}