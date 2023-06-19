package promise

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New[bool](func(ctx context.Context) (bool, error) {
		return false, nil
	})
	assert.NotNil(t, p)
}
