package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProprietyError_Error(t *testing.T) {
	e := proprietyError("message")
	assert.Equal(t, "message", e.Error())
}
