package promise

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMap(t *testing.T) {
	t.Run("against empty map", func(t *testing.T) {
		assert.PanicsWithValue(t, unexpectedCaseError("empty map"), func() {
			NewMap[any](map[any]any{})
		})
	})

	t.Run("against single-element map", func(t *testing.T) {
		m := NewMap[any](map[any]any{
			"my_promise": New(func(ctx context.Context) (any, error) {
				return "OK", nil
			}),
		})
		assert.IsType(t, &futureMap[any]{}, m)
	})

	t.Run("against multi-element map", func(t *testing.T) {
		m := NewMap[any](map[any]any{
			"my_promise": New(func(ctx context.Context) (any, error) {
				return "OK", nil
			}),
			1: New(func(ctx context.Context) (any, error) {
				return true, nil
			}),
		})
		assert.IsType(t, &futureMap[any]{}, m)
	})
}
