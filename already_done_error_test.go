package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlreadyDoneError_Error(t *testing.T) {
	e := alreadyDoneError("message")
	assert.Equal(t, "message", e.Error())
}
