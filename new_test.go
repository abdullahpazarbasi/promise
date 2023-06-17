package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New[bool](func() (bool, error) {
		return false, nil
	})
	assert.NotNil(t, p)
}
