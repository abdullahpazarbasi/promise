package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimedOutError_Error(t *testing.T) {
	e := timedOutError("message")
	assert.Equal(t, "message", e.Error())
}
