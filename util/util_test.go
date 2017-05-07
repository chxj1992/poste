package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToAddr(t *testing.T) {
	assert.Equal(t, "127.0.0.1:12345", ToAddr("127.0.0.1", 12345))
}

func TestRandom(t *testing.T) {
	assert.Contains(t, []string{"a", "b", "c"}, Random([]string{"a", "b", "c"}))
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", Md5("hello world"))
}

func TestBase64Encode(t *testing.T) {
	assert.Equal(t, "aGVsbG8gd29ybGQ=", Base64Encode("hello world"))
}

func TestBase64Decode(t *testing.T) {
	assert.Equal(t, "hello world", Base64Decode("aGVsbG8gd29ybGQ="))
}

func TestRandStr(t *testing.T) {
	assert.Equal(t, 8, len(RandStr(8)))
}
