package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	v := 42
	p := Ptr(v)

	assert.NotNil(t, p)
	assert.Equal(t, 42, *p)

	v = 100
	assert.Equal(t, 42, *p)
}
