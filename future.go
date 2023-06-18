package promise

import (
	"context"
	"sync"
	"time"
)

type Future[T any] interface {
	TimeOutLimit(timeOutLimit time.Duration) Future[T]
	OnResolved(onResolved func(T)) Future[T]
	OnRejected(onRejected func(error)) Future[T]
	OnCanceled(onCanceled func()) Future[T]
	OnTimedOut(onTimedOut func()) Future[T]
	Finally(finally func(event)) Future[T]
	Commit() Progress[T]
	Await() (T, error)
	commit() Progress[T]
	setKey(key interface{}) Future[T]
	setFulfilmentChannel(fulfilmentChannel chan Output[T]) Future[T]
	setContext(cancelContext context.Context, cancel context.CancelFunc) Future[T]
}

type future[T any] struct {
	async                    func() (T, error)
	timeOutLimit             time.Duration
	doOnResolved             func(T)
	doOnRejected             func(error)
	doOnCanceled             func()
	doOnTimedOut             func()
	doFinally                func(event)
	fulfilmentChannel        chan Output[T]
	cancelContext            context.Context
	cancel                   context.CancelFunc
	key                      interface{}
	timeOutLimitSetOnce      sync.Once
	doOnResolvedSetOnce      sync.Once
	doOnRejectedSetOnce      sync.Once
	doOnCanceledSetOnce      sync.Once
	doOnTimedOutSetOnce      sync.Once
	doFinallySetOnce         sync.Once
	fulfilmentChannelSetOnce sync.Once
	cancelContextSetOnce     sync.Once
	committedOnce            sync.Once
}

func (f *future[T]) TimeOutLimit(timeOutLimit time.Duration) Future[T] {
	f.timeOutLimitSetOnce.Do(func() {
		f.timeOutLimit = timeOutLimit
	})

	return f
}

func (f *future[T]) OnResolved(doOnResolved func(T)) Future[T] {
	f.doOnResolvedSetOnce.Do(func() {
		f.doOnResolved = doOnResolved
	})

	return f
}

func (f *future[T]) OnRejected(doOnRejected func(error)) Future[T] {
	f.doOnRejectedSetOnce.Do(func() {
		f.doOnRejected = doOnRejected
	})

	return f
}

func (f *future[T]) OnCanceled(doOnCanceled func()) Future[T] {
	f.doOnCanceledSetOnce.Do(func() {
		f.doOnCanceled = doOnCanceled
	})

	return f
}

func (f *future[T]) OnTimedOut(onTimedOut func()) Future[T] {
	f.doOnTimedOutSetOnce.Do(func() {
		f.doOnTimedOut = onTimedOut
	})

	return f
}

func (f *future[T]) Finally(doFinally func(event)) Future[T] {
	f.doFinallySetOnce.Do(func() {
		f.doFinally = doFinally
	})

	return f
}

func (f *future[T]) Commit() Progress[T] {
	return f.commit()
}

func (f *future[T]) Await() (T, error) {
	return f.commit().await()
}

func (f *future[T]) commit() Progress[T] {
	var pro *progress[T]
	f.committedOnce.Do(func() {
		if f.doOnTimedOut != nil && f.timeOutLimit == 0 {
			panic(proprietyError("on-timed-out is determined although time-out limit is not declared"))
		}
		f.fulfilmentChannelSetOnce.Do(func() {
			f.fulfilmentChannel = make(chan Output[T], 1)
		})
		f.cancelContextSetOnce.Do(func() {
			ctx := context.Background()
			if f.timeOutLimit == 0 {
				f.cancelContext, f.cancel = context.WithCancel(ctx)
			} else {
				f.cancelContext, f.cancel = context.WithTimeout(ctx, f.timeOutLimit)
			}
		})
		pro = &progress[T]{
			doOnResolved:      f.doOnResolved,
			doOnRejected:      f.doOnRejected,
			doOnCanceled:      f.doOnCanceled,
			doOnTimedOut:      f.doOnTimedOut,
			doFinally:         f.doFinally,
			fulfilmentChannel: f.fulfilmentChannel,
			context:           f.cancelContext,
			cancel:            f.cancel,
			key:               f.key,
			doneOnce:          sync.Once{},
		}

		go func() {
			defer pro.handleProbablePanic()

			pay, err := f.async()
			if err != nil {
				pro.reject(err)

				return
			}
			pro.resolve(pay)
		}()
	})

	return pro
}

func (f *future[T]) setKey(key interface{}) Future[T] {
	f.key = key

	return f
}

func (f *future[T]) setFulfilmentChannel(fulfilmentChannel chan Output[T]) Future[T] {
	f.fulfilmentChannelSetOnce.Do(func() {
		f.fulfilmentChannel = fulfilmentChannel
	})

	return f
}

func (f *future[T]) setContext(cancelContext context.Context, cancel context.CancelFunc) Future[T] {
	f.cancelContextSetOnce.Do(func() {
		f.cancelContext = cancelContext
		f.cancel = cancel
	})

	return f
}

func (f *future[T]) getCancelContext() context.Context {
	return f.cancelContext
}

func (f *future[T]) getCancelFunction() context.CancelFunc {
	return f.cancel
}
