package promise

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFuture_TimeOutLimit(t *testing.T) {
	actualPromise := New[string](func() (string, error) {
		return "", nil
	}).TimeOutLimit(200 * time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, actualPromise.(*future[string]).timeOutLimit)
}

func TestFuture_OnResolved(t *testing.T) {
	var jar string
	New(func() (string, error) {
		return "OK", nil
	}).OnResolved(func(s string) {
		jar = s
	}).Await()
	assert.Equal(t, "OK", jar)
}

func TestFuture_OnRejected(t *testing.T) {
	var jar error
	New[string](func() (string, error) {
		return "", fmt.Errorf("failed")
	}).OnRejected(func(err error) {
		jar = err
	}).Await()
	assert.Equal(t, "failed", jar.Error())
}

func TestFuture_Finally(t *testing.T) {
	t.Run("against resolution", func(t *testing.T) {
		var hit bool
		New[string](func() (string, error) {
			return "OK", nil
		}).Finally(func(e event) {
			hit = true
		}).Await()
		assert.True(t, hit)
	})

	t.Run("against rejection", func(t *testing.T) {
		var hit bool
		New(func() (string, error) {
			return "", fmt.Errorf("failed")
		}).Finally(func(e event) {
			hit = true
		}).Await()
		assert.True(t, hit)
	})
}

func TestFuture_OnCanceled(t *testing.T) {
	var resolved, rejected, canceled, timedOut, finished bool
	p := New(func() (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		10 * time.Second,
	).OnResolved(func(s string) {
		resolved = true
	}).OnRejected(func(err error) {
		rejected = true
	}).OnCanceled(func() {
		canceled = true
	}).OnTimedOut(func() {
		timedOut = true
	}).Finally(func(e event) {
		finished = true
	}).Commit()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Cancel()
	}()
	p.Await()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.True(t, canceled)
	assert.False(t, timedOut)
	assert.True(t, finished)
}

func TestFuture_OnTimedOut(t *testing.T) {
	var resolved, rejected, canceled, timedOut, finished bool
	New[string](func() (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		200 * time.Millisecond,
	).OnResolved(func(s string) {
		resolved = true
	}).OnRejected(func(err error) {
		rejected = true
	}).OnCanceled(func() {
		canceled = true
	}).OnTimedOut(func() {
		timedOut = true
	}).Finally(func(e event) {
		finished = true
	}).Await()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.False(t, canceled)
	assert.True(t, timedOut)
	assert.True(t, finished)
}

func TestFuture_Commit(t *testing.T) {
	t.Run("against fulfilment", func(t *testing.T) {
		actualResult, err := New(func() (string, error) {
			return "OK", nil
		}).Await()
		if assert.NoError(t, err) {
			assert.Equal(t, "OK", actualResult)
		}
	})

	t.Run("against failure", func(t *testing.T) {
		actualResult, err := New[string](func() (string, error) {
			return "a", fmt.Errorf("failed")
		}).Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "failed", err.Error())
		}
	})

	t.Run("against panic", func(t *testing.T) {
		actualResult, err := New(func() (string, error) {
			panic("aaaaaaaaaa")
		}).Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "aaaaaaaaaa", err.Error())
		}
	})
}
