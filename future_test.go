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
	assert.Equal(t, 200*time.Millisecond, actualPromise.(*Future[string]).timeOutLimit)
}

func TestFuture_OnResolved(t *testing.T) {
	var jar string
	New(func() (string, error) {
		return "OK", nil
	}).OnResolved(func(s string) {
		jar = s
	}).Commit().Await()
	assert.Equal(t, "OK", jar)
}

func TestFuture_OnRejected(t *testing.T) {
	var jar error
	New[string](func() (string, error) {
		return "", fmt.Errorf("failed")
	}).OnRejected(func(err error) {
		jar = err
	}).Commit().Await()
	assert.Equal(t, "failed", jar.Error())
}

func TestFuture_OnCompleted(t *testing.T) {
	t.Run("against resolution", func(t *testing.T) {
		var hit bool
		New[string](func() (string, error) {
			return "OK", nil
		}).OnCompleted(func() {
			hit = true
		}).Commit().Await()
		assert.True(t, hit)
	})

	t.Run("against rejection", func(t *testing.T) {
		var hit bool
		New(func() (string, error) {
			return "", fmt.Errorf("failed")
		}).OnCompleted(func() {
			hit = true
		}).Commit().Await()
		assert.True(t, hit)
	})
}

func TestFuture_OnCanceled(t *testing.T) {
	var resolved, rejected, completed, canceled, timedOut bool
	p := New(func() (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		10 * time.Second,
	).OnResolved(func(s string) {
		resolved = true
	}).OnRejected(func(err error) {
		rejected = true
	}).OnCompleted(func() {
		completed = true
	}).OnCanceled(func() {
		canceled = true
	}).OnTimedOut(func() {
		timedOut = true
	}).Commit()
	go func() {
		time.Sleep(200 * time.Millisecond)
		p.Cancel()
	}()
	p.Await()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.False(t, completed)
	assert.True(t, canceled)
	assert.False(t, timedOut)
}

func TestFuture_OnTimedOut(t *testing.T) {
	var resolved, rejected, completed, canceled, timedOut bool
	p := New[string](func() (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		200 * time.Millisecond,
	).OnResolved(func(s string) {
		resolved = true
	}).OnRejected(func(err error) {
		rejected = true
	}).OnCompleted(func() {
		completed = true
	}).OnCanceled(func() {
		canceled = true
	}).OnTimedOut(func() {
		timedOut = true
	}).Commit()
	p.Await()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.False(t, completed)
	assert.False(t, canceled)
	assert.True(t, timedOut)
}

func TestFuture_Commit(t *testing.T) {
	t.Run("against fulfilment", func(t *testing.T) {
		actualResult, err := New(func() (string, error) {
			return "OK", nil
		}).Commit().Await()
		if assert.NoError(t, err) {
			assert.Equal(t, "OK", actualResult)
		}
	})

	t.Run("against failure", func(t *testing.T) {
		actualResult, err := New[string](func() (string, error) {
			return "a", fmt.Errorf("failed")
		}).Commit().Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "failed", err.Error())
		}
	})

	t.Run("against panic", func(t *testing.T) {
		actualResult, err := New(func() (string, error) {
			panic("aaaaaaaaaa")
		}).Commit().Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "aaaaaaaaaa", err.Error())
		}
	})
}
