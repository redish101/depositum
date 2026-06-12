package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEcho(t *testing.T) {
	e := NewEcho()
	assert.NotNil(t, e)
}
