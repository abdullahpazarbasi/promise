package promise

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProgress_Cancel(t *testing.T) {
	t.Run("against parallel call", func(t *testing.T) {
		p := New[bool]().Commit(func() (bool, error) {
			time.Sleep(500 * time.Millisecond)

			return true, nil
		})
		go func() {
			time.Sleep(200 * time.Millisecond)
			p.Cancel()
		}()
		actualResult, err := p.Await()
		assert.False(t, actualResult)
		if assert.Error(t, err) {
			assert.IsType(t, context.Canceled, err)
		}
	})

	t.Run("against paid-off promise", func(t *testing.T) {
		p := New[bool]().Commit(func() (bool, error) {
			time.Sleep(500 * time.Millisecond)

			return true, nil
		})
		actualResult, err := p.Await()
		p.Cancel()
		assert.True(t, actualResult)
		assert.NoError(t, err)
	})
}

func TestProgress_Await(t *testing.T) {
	t.Run("against no deadline", func(t *testing.T) {
		p := New[bool]()
		actualResult, err := p.Commit(func() (bool, error) {
			time.Sleep(500 * time.Millisecond)

			return true, nil
		}).Await()
		if assert.NoError(t, err) {
			assert.True(t, actualResult)
		}
	})

	t.Run("against early deadline", func(t *testing.T) {
		p := New[bool]().TimeOutLimit(100 * time.Millisecond)
		actualResult, err := p.Commit(func() (bool, error) {
			time.Sleep(500 * time.Millisecond)

			return true, nil
		}).Await()
		assert.False(t, actualResult)
		if assert.Error(t, err) {
			assert.IsType(t, context.DeadlineExceeded, err)
		}
	})
}
