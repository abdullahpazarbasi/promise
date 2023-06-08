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
	Commit(async func() (T, error)) Promise[T]
	Cancel() Promise[T]
	Await() (T, error)
}

func New[T any]() Promise[T] {
	return &Future[T]{
		timeOutLimit:      0,
		fulfilmentChannel: make(chan struct{}),
		committedOnce:     sync.Once{},
	}
}
