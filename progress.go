package promise

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Progress[T any] struct {
	doOnResolved      func(T)
	doOnRejected      func(error)
	doOnCompleted     func()
	doOnCanceled      func()
	doOnTimedOut      func()
	doFinally         func(event)
	fulfilmentChannel chan struct {
		out *T
		err error
	}
	context  context.Context
	cancel   context.CancelFunc
	doneOnce sync.Once
}

func (p *Progress[T]) TimeOutLimit(duration time.Duration) Promise[T] {
	panic(proprietyError("time-out limit can not be determined in progress"))
}

func (p *Progress[T]) OnResolved(onResolved func(T)) Promise[T] {
	panic(proprietyError("on-resolved can not be determined in progress"))
}

func (p *Progress[T]) OnRejected(onRejected func(error)) Promise[T] {
	panic(proprietyError("on-rejected can not be determined in progress"))
}

func (p *Progress[T]) OnCompleted(onCompleted func()) Promise[T] {
	panic(proprietyError("on-completed can not be determined in progress"))
}

func (p *Progress[T]) OnCanceled(onCanceled func()) Promise[T] {
	panic(proprietyError("on-canceled can not be determined in progress"))
}

func (p *Progress[T]) OnTimedOut(onTimedOut func()) Promise[T] {
	panic(proprietyError("on-timed-out can not be determined in progress"))
}

func (p *Progress[T]) Commit() Promise[T] {
	return p
}

func (p *Progress[T]) Cancel() {
	p.cancel()
}

func (p *Progress[T]) Await() (T, error) {
	select {
	case <-p.context.Done():
		return p.abandon(p.context.Err())
	case r := <-p.fulfilmentChannel:
		return p.fulfil(r)
	}
}

func (p *Progress[T]) getFulfilmentChannel() chan struct {
	out *T
	err error
} {
	return p.fulfilmentChannel
}

func (p *Progress[T]) getContext() context.Context {
	return p.context
}

func (p *Progress[T]) fulfil(result struct {
	out *T
	err error
}) (T, error) {
	defer p.cancel()

	if result.out == nil {
		var defaultT T

		return defaultT, result.err
	}

	return *result.out, result.err
}

func (p *Progress[T]) abandon(err error) (T, error) {
	var defaultT T
	switch err {
	case context.Canceled:
		if p.doOnCanceled != nil {
			p.doOnCanceled()
		}
		if p.doFinally != nil {
			p.doFinally(EventCanceled)
		}
	case context.DeadlineExceeded:
		if p.doOnTimedOut != nil {
			p.doOnTimedOut()
		}
		if p.doFinally != nil {
			p.doFinally(EventTimedOut)
		}
	default:
		panic(unexpectedCaseError("unexpected error type"))
	}

	return defaultT, err
}

func (p *Progress[T]) handleProbablePanic() {
	err := recover()
	if err == nil {
		return
	}

	switch e := err.(type) {
	case error:
		p.reject(e)
	default:
		p.reject(fmt.Errorf("%+v", e))
	}
}

func (p *Progress[T]) resolve(out T) {
	p.doneOnce.Do(func() {
		p.fulfilmentChannel <- struct {
			out *T
			err error
		}{
			out: &out,
			err: nil,
		}
		close(p.fulfilmentChannel)
		if p.doOnResolved != nil {
			p.doOnResolved(out)
		}
		if p.doOnCompleted != nil {
			p.doOnCompleted()
		}
		if p.doFinally != nil {
			p.doFinally(EventResolved)
		}
	})
}

func (p *Progress[T]) reject(err error) {
	p.doneOnce.Do(func() {
		p.fulfilmentChannel <- struct {
			out *T
			err error
		}{
			out: nil,
			err: err,
		}
		close(p.fulfilmentChannel)
		if p.doOnRejected != nil {
			p.doOnRejected(err)
		}
		if p.doOnCompleted != nil {
			p.doOnCompleted()
		}
		if p.doFinally != nil {
			p.doFinally(EventRejected)
		}
	})
}
