package promise

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMap(t *testing.T) {
	t.Run("against empty map", func(t *testing.T) {
		assert.PanicsWithValue(t, unexpectedCaseError("empty map"), func() {
			NewMap(map[any]any{}, 0)
		})
	})

	t.Run("against single-element map", func(t *testing.T) {
		m := NewMap(map[any]any{
			"my_promise": New(func() (any, error) {
				return "OK", nil
			}),
		}, 0)
		assert.IsType(t, &promiseMap{}, m)
	})

	t.Run("against multi-element map", func(t *testing.T) {
		m := NewMap(map[any]any{
			"my_promise": New(func() (any, error) {
				return "OK", nil
			}),
			1: New(func() (any, error) {
				return true, nil
			}),
		}, 0)
		assert.IsType(t, &promiseMap{}, m)
	})
}

func TestConcreteMap_AssureCommitments(t *testing.T) {
	//
}

func TestConcreteMap_Cancel(t *testing.T) {
	//
}

func TestConcreteMap_Await(t *testing.T) {
	t.Run("against 2 resolvable futures", func(t *testing.T) {
		m := NewMap(
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
		actualResultMap := m.AssureCommitments().Await()
		assert.Equal(
			t,
			map[interface{}]struct {
				out interface{}
				err error
			}{
				"my_promise": {
					out: "OK",
					err: nil,
				},
				1: {
					out: true,
					err: nil,
				},
			},
			*actualResultMap,
		)
	})

	t.Run("against 2 rejectable futures", func(t *testing.T) {
		m := NewMap(
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
		actualResultMap := m.AssureCommitments().Await()
		assert.Equal(
			t,
			map[interface{}]struct {
				out interface{}
				err error
			}{
				"my_promise": {
					out: nil,
					err: fmt.Errorf("oops"),
				},
				1: {
					out: nil,
					err: fmt.Errorf("oops"),
				},
			},
			*actualResultMap,
		)
	})

	t.Run("against 2 futures, one of them taking too long", func(t *testing.T) {
		m := NewMap(
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
		actualResultMap := m.AssureCommitments().Await()
		assert.Equal(
			t,
			map[interface{}]struct {
				out interface{}
				err error
			}{
				"my_promise": {
					out: nil,
					err: context.DeadlineExceeded,
				},
				1: {
					out: true,
					err: nil,
				},
			},
			*actualResultMap,
		)
	})
}

func TestConcreteMap_Race(t *testing.T) {
	//
}
