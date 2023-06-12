package promise

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	//
}

func TestConcreteMap_Race(t *testing.T) {
	//
}
