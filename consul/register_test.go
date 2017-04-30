package consul

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert.Nil(t, Register("test", "127.0.0.1", 12345, nil))
}