package data

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestDataEncode(t *testing.T) {
	data := Data{Target:"1", Message:"hello world"}
	bytes := data.Marshal()
	assert.Equal(t, "{\"target\":\"1\",\"message\":\"hello world\"}", string(bytes))
}

func TestDataDecode(t *testing.T) {
	bytes := []byte("{\"target\":\"1\",\"message\":\"hello world\"}")
	var data Data
	json.Unmarshal(bytes, &data)
	assert.Equal(t, Data{Target:"1", Message:"hello world"}, data)
}
