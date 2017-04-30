package dispather

import (
	"encoding/json"
	"net/http"
)

func Mailmen(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(mailmen)
	w.Write(bytes)
}