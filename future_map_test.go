package promise

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFutureMap_Commit(t *testing.T) {
	// todo:
}

func TestFutureMap_Await(t *testing.T) {
	t.Run("against 2 resolvable futures", func(t *testing.T) {
		m := NewMap[any](
			map[any]any{
				"my_promise": New(func() (any, error) {
					time.Sleep(300 * time.Millisecond)

					return "OK", nil
				}),
				1: New(func() (any, error) {
					time.Sleep(300 * time.Millisecond)

					return true, nil
				}),
			},
			400*time.Millisecond,
		)
		actualResultMap := m.Await()
		assert.Equal(
			t,
			map[interface{}]Output[any]{
				"my_promise": newOutput[any]("OK", nil),
				1:            newOutput[any](true, nil),
			},
			*actualResultMap,
		)
	})

	t.Run("against 2 rejectable futures", func(t *testing.T) {
		m := NewMap[any](
			map[any]any{
				"my_promise": New(func() (any, error) {
					time.Sleep(300 * time.Millisecond)

					return "", fmt.Errorf("oops")
				}),
				1: New(func() (any, error) {
					time.Sleep(300 * time.Millisecond)

					return false, fmt.Errorf("oops")
				}),
			},
			400*time.Millisecond,
		)
		actualResultMap := m.Await()
		assert.Equal(
			t,
			map[interface{}]Output[any]{
				"my_promise": newOutput[any](nil, fmt.Errorf("oops")),
				1:            newOutput[any](nil, fmt.Errorf("oops")),
			},
			*actualResultMap,
		)
	})

	t.Run("against 2 futures, one of them taking too long", func(t *testing.T) {
		m := NewMap[any](
			map[any]any{
				"my_promise": New(func() (any, error) {
					time.Sleep(500 * time.Millisecond)

					return "OK", nil
				}),
				1: New(func() (any, error) {
					time.Sleep(100 * time.Millisecond)

					return true, nil
				}),
			},
			300*time.Millisecond,
		)
		actualResultMap := m.Await()
		assert.Equal(
			t,
			map[interface{}]Output[any]{
				"my_promise": newOutput[any](nil, context.DeadlineExceeded),
				1:            newOutput[any](nil, context.DeadlineExceeded),
			},
			*actualResultMap,
		)
	})
}

func TestFutureMap_Race(t *testing.T) {
	// todo:
}
