package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New[bool]()
	assert.NotNil(t, p)
}
