package promise

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

	t.Run("against 3 futures, one of them taking too long", func(t *testing.T) {
		m := NewMap[any](
			map[any]any{
				"my_promise": New(func() (any, error) {
					time.Sleep(100 * time.Millisecond)

					return "OK", nil
				}),
				1: New(func() (any, error) {
					time.Sleep(500 * time.Millisecond)

					return 1, nil
				}),
				true: New(func() (any, error) {
					time.Sleep(200 * time.Millisecond)

					return true, nil
				}),
			},
			300*time.Millisecond,
		)
		actualResultMap := m.Await()
		assert.Equal(
			t,
			map[interface{}]Output[any]{
				1:            newOutput[any](nil, timedOutError("timed-out")),
				true:         newOutput[any](true, nil),
				"my_promise": newOutput[any]("OK", nil),
			},
			*actualResultMap,
		)
	})
}

func TestFutureMap_Race(t *testing.T) {
	t.Run("against 2 resolvable futures", func(t *testing.T) {
		m := NewMap[any](
			map[any]any{
				"my_promise": New(func() (any, error) {
					time.Sleep(300 * time.Millisecond)

					return "OK", nil
				}),
				1: New(func() (any, error) {
					time.Sleep(200 * time.Millisecond)

					return true, nil
				}),
			},
			400*time.Millisecond,
		)
		actualKey, actualPayload, actualError := m.Race()
		assert.Equal(t, 1, actualKey)
		assert.Equal(t, true, actualPayload)
		assert.Equal(t, nil, actualError)
	})
}
