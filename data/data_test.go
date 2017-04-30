package data

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestDataEncode(t *testing.T) {
	data := Data{Target:"1", ServerType:"ws", Message:"hello world"}
	bytes, _ := json.Marshal(data)
	assert.Equal(t, "{\"target\":\"1\",\"type\":\"ws\",\"message\":\"hello world\"}", string(bytes))
}

func TestDataDecode(t *testing.T) {
	bytes := []byte("{\"target\":\"1\",\"type\":\"ws\",\"message\":\"hello world\"}")
	var data Data
	json.Unmarshal(bytes, &data)
	assert.Equal(t, Data{Target:"1", ServerType:"ws", Message:"hello world"}, data)
}
