package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanceledError_Error(t *testing.T) {
	e := canceledError("message")
	assert.Equal(t, "message", e.Error())
}
