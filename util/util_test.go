package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToAddr(t *testing.T) {
	assert.Equal(t, "127.0.0.1:12345", ToAddr("127.0.0.1", 12345))
}
