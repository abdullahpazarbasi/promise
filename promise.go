package promise

import (
	"sync"
	"time"
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
}

func New[T any](async func() (T, error)) Promise[T] {
	return &Future[T]{
		async:             async,
		timeOutLimit:      0,
		fulfilmentChannel: make(chan struct{}),
		committedOnce:     sync.Once{},
	}
}
