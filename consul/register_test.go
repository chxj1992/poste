package consul

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert.Nil(t, Register("test", "127.0.0.1", 12345, nil))
}

func TestServiceId(t *testing.T) {
	assert.Equal(t, "test_8e4c009607bc3b56b19a05969fdb8d9a", ServiceId("test", "127.0.0.1", 12345))
}