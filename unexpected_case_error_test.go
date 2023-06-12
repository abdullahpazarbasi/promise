package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnexpectedCaseError_Error(t *testing.T) {
	e := unexpectedCaseError("message")
	assert.Equal(t, "message", e.Error())
}
