package promise

import (
	"context"
	"sync"
	"time"
)

type event int

const (
	EventResolved event = iota
	EventRejected
	EventCanceled
	EventTimedOut
)

type Promise[T any] interface {
	TimeOutLimit(duration time.Duration) Promise[T]
	OnResolved(onResolved func(T)) Promise[T]
	OnRejected(onRejected func(error)) Promise[T]
	OnCompleted(onCompleted func()) Promise[T]
	OnCanceled(onCanceled func()) Promise[T]
	OnTimedOut(onTimedOut func()) Promise[T]
	Commit() Promise[T]
	Cancel()
	Await() (T, error)
	getFulfilmentChannel() chan struct {
		out *T
		err error
	}
	getContext() context.Context
	fulfil(result struct {
		out *T
		err error
	}) (T, error)
	abandon(err error) (T, error)
}

func New[T any](async func() (T, error)) Promise[T] {
	return &Future[T]{
		async:         async,
		timeOutLimit:  0,
		committedOnce: sync.Once{},
	}
}
