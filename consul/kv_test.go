package consul

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func setUp() {
	KVClear()
	KVSet("hello", "world")
	KVSet("name/a", "Tom")
	KVSet("name/b", "Andy")
}

func TestGet(t *testing.T) {
	setUp()

	assert.Equal(t, "world", KVGet("hello"))
}

func TestSet(t *testing.T) {
	setUp()

	assert.True(t, KVSet("hello", "Tom"))
	assert.Equal(t, "Tom", KVGet("hello"))
}

func TestDelete(t *testing.T) {
	setUp()

	assert.True(t, KVDelete("hello"))
	assert.Equal(t, "", KVGet("hello"))
}

func TestClear(t *testing.T) {
	setUp()

	assert.True(t, KVClear())
	assert.Equal(t, "", KVGet("hello"))
}

func TestValues(t *testing.T) {
	setUp()

	assert.Equal(t, []string{"Tom", "Andy"}, KVValues("name"))
}
