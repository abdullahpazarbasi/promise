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

func (p *future[T]) TimeOutLimit(timeOutLimit time.Duration) Future[T] {
	p.timeOutLimitSetOnce.Do(func() {
		p.timeOutLimit = timeOutLimit
	})

	return p
}

func (p *future[T]) OnResolved(doOnResolved func(T)) Future[T] {
	p.doOnResolvedSetOnce.Do(func() {
		p.doOnResolved = doOnResolved
	})

	return p
}

func (p *future[T]) OnRejected(doOnRejected func(error)) Future[T] {
	p.doOnRejectedSetOnce.Do(func() {
		p.doOnRejected = doOnRejected
	})

	return p
}

func (p *future[T]) OnCanceled(doOnCanceled func()) Future[T] {
	p.doOnCanceledSetOnce.Do(func() {
		p.doOnCanceled = doOnCanceled
	})

	return p
}

func (p *future[T]) OnTimedOut(onTimedOut func()) Future[T] {
	p.doOnTimedOutSetOnce.Do(func() {
		p.doOnTimedOut = onTimedOut
	})

	return p
}

func (p *future[T]) Finally(doFinally func(event)) Future[T] {
	p.doFinallySetOnce.Do(func() {
		p.doFinally = doFinally
	})

	return p
}

func (p *future[T]) Commit() Progress[T] {
	return p.commit()
}

func (p *future[T]) Await() (T, error) {
	return p.commit().await()
}

func (p *future[T]) commit() Progress[T] {
	var ps *progress[T]
	p.committedOnce.Do(func() {
		if p.doOnTimedOut != nil && p.timeOutLimit == 0 {
			panic(proprietyError("on-timed-out is determined although time-out limit is not declared"))
		}
		p.fulfilmentChannelSetOnce.Do(func() {
			p.fulfilmentChannel = make(chan Output[T], 1)
		})
		p.cancelContextSetOnce.Do(func() {
			ctx := context.Background()
			if p.timeOutLimit == 0 {
				p.cancelContext, p.cancel = context.WithCancel(ctx)
			} else {
				p.cancelContext, p.cancel = context.WithTimeout(ctx, p.timeOutLimit)
			}
		})
		ps = &progress[T]{
			doOnResolved:      p.doOnResolved,
			doOnRejected:      p.doOnRejected,
			doOnCanceled:      p.doOnCanceled,
			doOnTimedOut:      p.doOnTimedOut,
			doFinally:         p.doFinally,
			fulfilmentChannel: p.fulfilmentChannel,
			context:           p.cancelContext,
			cancel:            p.cancel,
			key:               p.key,
			doneOnce:          sync.Once{},
		}

		go func() {
			defer ps.handleProbablePanic()

			out, err := p.async()
			if err != nil {
				ps.reject(err)

				return
			}
			ps.resolve(out)
		}()
	})

	return ps
}

func (p *future[T]) setKey(key interface{}) Future[T] {
	p.key = key

	return p
}

func (p *future[T]) setFulfilmentChannel(fulfilmentChannel chan Output[T]) Future[T] {
	p.fulfilmentChannelSetOnce.Do(func() {
		p.fulfilmentChannel = fulfilmentChannel
	})

	return p
}

func (p *future[T]) setContext(cancelContext context.Context, cancel context.CancelFunc) Future[T] {
	p.cancelContextSetOnce.Do(func() {
		p.cancelContext = cancelContext
		p.cancel = cancel
	})

	return p
}

func (p *future[T]) getCancelContext() context.Context {
	return p.cancelContext
}

func (p *future[T]) getCancelFunction() context.CancelFunc {
	return p.cancel
}
