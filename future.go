package promise

import (
	"context"
	"sync"
	"time"
)

type Future[T any] struct {
	async                func() (T, error)
	timeOutLimit         time.Duration
	doOnResolved         func(T)
	doOnRejected         func(error)
	doOnCompleted        func()
	doOnCanceled         func()
	doOnTimedOut         func()
	doFinally            func(event)
	timeOutLimitSetOnce  sync.Once
	doOnResolvedSetOnce  sync.Once
	doOnRejectedSetOnce  sync.Once
	doOnCompletedSetOnce sync.Once
	doOnCanceledSetOnce  sync.Once
	doOnTimedOutSetOnce  sync.Once
	committedOnce        sync.Once
}

func (p *Future[T]) TimeOutLimit(duration time.Duration) Promise[T] {
	p.timeOutLimitSetOnce.Do(func() {
		p.timeOutLimit = duration
	})

	return p
}

func (p *Future[T]) OnResolved(doOnResolved func(T)) Promise[T] {
	p.doOnResolvedSetOnce.Do(func() {
		p.doOnResolved = doOnResolved
	})

	return p
}

func (p *Future[T]) OnRejected(doOnRejected func(error)) Promise[T] {
	p.doOnRejectedSetOnce.Do(func() {
		p.doOnRejected = doOnRejected
	})

	return p
}

func (p *Future[T]) OnCompleted(doOnCompleted func()) Promise[T] {
	p.doOnCompletedSetOnce.Do(func() {
		p.doOnCompleted = doOnCompleted
	})

	return p
}

func (p *Future[T]) OnCanceled(doOnCanceled func()) Promise[T] {
	p.doOnCanceledSetOnce.Do(func() {
		p.doOnCanceled = doOnCanceled
	})

	return p
}

func (p *Future[T]) OnTimedOut(onTimedOut func()) Promise[T] {
	p.doOnTimedOutSetOnce.Do(func() {
		p.doOnTimedOut = onTimedOut
	})

	return p
}

func (p *Future[T]) Commit() Promise[T] {
	var progress *Progress[T]
	p.committedOnce.Do(func() {
		if p.doOnTimedOut != nil && p.timeOutLimit == 0 {
			panic(proprietyError("on-timed-out is determined although time-out limit is not declared"))
		}
		progress = &Progress[T]{
			doOnResolved:  p.doOnResolved,
			doOnRejected:  p.doOnRejected,
			doOnCompleted: p.doOnCompleted,
			doOnCanceled:  p.doOnCanceled,
			doOnTimedOut:  p.doOnTimedOut,
			doFinally:     p.doFinally,
			fulfilmentChannel: make(chan struct {
				out *T
				err error
			}),
			doneOnce: sync.Once{},
		}
		ctx := context.Background()
		if p.timeOutLimit == 0 {
			progress.context, progress.cancel = context.WithCancel(ctx)
		} else {
			progress.context, progress.cancel = context.WithTimeout(ctx, p.timeOutLimit)
		}

		go func() {
			defer progress.handleProbablePanic()

			out, err := p.async()
			if err != nil {
				progress.reject(err)

				return
			}
			progress.resolve(out)
		}()
	})

	return progress
}

func (p *Future[T]) Cancel() {
	panic(proprietyError("a promise which is not committed can not be canceled"))
}

func (p *Future[T]) Await() (T, error) {
	panic(proprietyError("a promise which is not committed can not be awaited"))
}

func (p *Future[T]) getFulfilmentChannel() chan struct {
	out *T
	err error
} {
	return nil
}

func (p *Future[T]) getContext() context.Context {
	return nil
}

func (p *Future[T]) fulfil(result struct {
	out *T
	err error
}) (T, error) {
	var defaultT T

	return defaultT, nil
}

func (p *Future[T]) abandon(err error) (T, error) {
	var defaultT T

	return defaultT, nil
}
