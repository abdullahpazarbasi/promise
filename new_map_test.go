package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMap(t *testing.T) {
	t.Run("against empty map", func(t *testing.T) {
		assert.PanicsWithValue(t, unexpectedCaseError("empty map"), func() {
			NewMap[any](map[any]any{}, 0)
		})
	})

	t.Run("against single-element map", func(t *testing.T) {
		m := NewMap[any](map[any]any{
			"my_promise": New(func() (any, error) {
				return "OK", nil
			}),
		}, 0)
		assert.IsType(t, &futureMap[any]{}, m)
	})

	t.Run("against multi-element map", func(t *testing.T) {
		m := NewMap[any](map[any]any{
			"my_promise": New(func() (any, error) {
				return "OK", nil
			}),
			1: New(func() (any, error) {
				return true, nil
			}),
		}, 0)
		assert.IsType(t, &futureMap[any]{}, m)
	})
}
