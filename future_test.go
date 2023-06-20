package promise

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestFuture_TimeOutLimit(t *testing.T) {
	actualPromise := New[string](func(ctx context.Context) (string, error) {
		return "", nil
	}).TimeOutLimit(200 * time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, actualPromise.(*future[string]).timeOutLimit)
}

func TestFuture_OnResolved(t *testing.T) {
	mu := sync.Mutex{}
	var jar string
	f := New(func(ctx context.Context) (string, error) {
		mu.Lock()
		jar = "OOPS"
		mu.Unlock()

		return "OK", nil
	}).OnResolved(func(s string) {
		mu.Lock()
		jar = s
		mu.Unlock()
	})
	f.Await()
	mu.Lock()
	assert.Equal(t, "OK", jar)
	mu.Unlock()
}

func TestFuture_OnRejected(t *testing.T) {
	mu := sync.Mutex{}
	var hit error
	f := New[string](func(ctx context.Context) (string, error) {
		mu.Lock()
		hit = fmt.Errorf("another fail")
		mu.Unlock()

		return "", fmt.Errorf("fail")
	}).OnRejected(func(err error) {
		mu.Lock()
		hit = err
		mu.Unlock()
	})
	f.Await()
	mu.Lock()
	assert.Equal(t, "fail", hit.Error())
	mu.Unlock()
}

func TestFuture_Finally(t *testing.T) {
	t.Run("against resolution", func(t *testing.T) {
		mu := sync.Mutex{}
		var hit bool
		New[string](func(ctx context.Context) (string, error) {
			return "OK", nil
		}).Finally(func(e event) {
			mu.Lock()
			hit = true
			mu.Unlock()
		}).Await()
		mu.Lock()
		assert.True(t, hit)
		mu.Unlock()
	})

	t.Run("against rejection", func(t *testing.T) {
		mu := sync.Mutex{}
		var hit bool
		New(func(ctx context.Context) (string, error) {
			return "", fmt.Errorf("failed")
		}).Finally(func(e event) {
			mu.Lock()
			hit = true
			mu.Unlock()
		}).Await()
		mu.Lock()
		assert.True(t, hit)
		mu.Unlock()
	})
}

func TestFuture_OnCanceled(t *testing.T) {
	mu := sync.Mutex{}
	var resolved, rejected, canceled, timedOut, finished bool
	p := New(func(ctx context.Context) (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		10 * time.Second,
	).OnResolved(func(s string) {
		mu.Lock()
		resolved = true
		mu.Unlock()
	}).OnRejected(func(err error) {
		mu.Lock()
		rejected = true
		mu.Unlock()
	}).OnCanceled(func() {
		mu.Lock()
		canceled = true
		mu.Unlock()
	}).OnTimedOut(func() {
		mu.Lock()
		timedOut = true
		mu.Unlock()
	}).Finally(func(e event) {
		mu.Lock()
		finished = true
		mu.Unlock()
	}).Commit()
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Cancel()
	}()
	p.Await()
	mu.Lock()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.True(t, canceled)
	assert.False(t, timedOut)
	assert.True(t, finished)
	mu.Unlock()
}

func TestFuture_OnTimedOut(t *testing.T) {
	mu := sync.Mutex{}
	var resolved, rejected, canceled, timedOut, finished bool
	New[string](func(ctx context.Context) (string, error) {
		time.Sleep(500 * time.Millisecond)

		return "OK", nil
	}).TimeOutLimit(
		200 * time.Millisecond,
	).OnResolved(func(s string) {
		mu.Lock()
		resolved = true
		mu.Unlock()
	}).OnRejected(func(err error) {
		mu.Lock()
		rejected = true
		mu.Unlock()
	}).OnCanceled(func() {
		mu.Lock()
		canceled = true
		mu.Unlock()
	}).OnTimedOut(func() {
		mu.Lock()
		timedOut = true
		mu.Unlock()
	}).Finally(func(e event) {
		mu.Lock()
		finished = true
		mu.Unlock()
	}).Await()
	mu.Lock()
	assert.False(t, resolved)
	assert.False(t, rejected)
	assert.False(t, canceled)
	assert.True(t, timedOut)
	assert.True(t, finished)
	mu.Unlock()
}

func TestFuture_Commit(t *testing.T) {
	t.Run("against fulfilment", func(t *testing.T) {
		actualResult, err := New(func(ctx context.Context) (string, error) {
			return "OK", nil
		}).Await()
		if assert.NoError(t, err) {
			assert.Equal(t, "OK", actualResult)
		}
	})

	t.Run("against failure", func(t *testing.T) {
		actualResult, err := New[string](func(ctx context.Context) (string, error) {
			return "a", fmt.Errorf("failed")
		}).Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "failed", err.Error())
		}
	})

	t.Run("against panic", func(t *testing.T) {
		actualResult, err := New(func(ctx context.Context) (string, error) {
			panic("aaaaaaaaaa")
		}).Await()
		assert.Equal(t, "", actualResult)
		if assert.Error(t, err) {
			assert.Equal(t, "aaaaaaaaaa", err.Error())
		}
	})
}
