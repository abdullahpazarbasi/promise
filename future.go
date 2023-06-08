package promise

import (
	"context"
	"sync"
	"time"
)

type Future[T any] struct {
	timeOutLimit         time.Duration
	doOnResolved         func(T)
	doOnRejected         func(error)
	doOnCompleted        func()
	doOnCanceled         func()
	doOnTimedOut         func()
	fulfilmentChannel    chan struct{}
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

func (p *Future[T]) Commit(async func() (T, error)) Promise[T] {
	var progress *Progress[T]
	p.committedOnce.Do(func() {
		if p.doOnTimedOut != nil && p.timeOutLimit == 0 {
			panic("on-timed-out is determined although time-out limit is not declared")
		}
		progress = &Progress[T]{
			doOnResolved:      p.doOnResolved,
			doOnRejected:      p.doOnRejected,
			doOnCompleted:     p.doOnCompleted,
			doOnCanceled:      p.doOnCanceled,
			doOnTimedOut:      p.doOnTimedOut,
			fulfilmentChannel: p.fulfilmentChannel,
			doneOnce:          sync.Once{},
		}
		ctx := context.Background()
		if p.timeOutLimit == 0 {
			progress.context, progress.cancel = context.WithCancel(ctx)
		} else {
			progress.context, progress.cancel = context.WithTimeout(ctx, p.timeOutLimit)
		}

		go func() {
			defer progress.handleProbablePanic()

			val, err := async()
			if err != nil {
				progress.reject(err)

				return
			}
			progress.resolve(val)
		}()
	})

	return progress
}

func (p *Future[T]) Cancel() Promise[T] {
	panic("a promise which is not committed can not be canceled")
}

func (p *Future[T]) Await() (T, error) {
	panic("a promise which is not committed can not be awaited")
}
