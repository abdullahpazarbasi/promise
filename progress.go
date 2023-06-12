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
	fulfilmentChannel chan struct{}
	outputResult      *T
	outputError       error
	context           context.Context
	cancel            context.CancelFunc
	doneOnce          sync.Once
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
	panic(alreadyDoneError("promise is already committed"))
}

func (p *Progress[T]) Cancel() {
	p.cancel()
}

func (p *Progress[T]) Await() (T, error) {
	select {
	case <-p.context.Done():
		var defaultT T
		err := p.context.Err()
		switch err {
		case context.Canceled:
			if p.doOnCanceled != nil {
				p.doOnCanceled()
			}
		case context.DeadlineExceeded:
			if p.doOnTimedOut != nil {
				p.doOnTimedOut()
			}
		default:
			panic(unexpectedCaseError("unexpected error type"))
		}

		return defaultT, err
	case <-p.fulfilmentChannel:
		defer p.cancel()

		if p.outputResult == nil {
			var defaultT T

			return defaultT, p.outputError
		}

		return *p.outputResult, p.outputError
	}
}

func (p *Progress[T]) resolve(val T) {
	p.doneOnce.Do(func() {
		p.outputResult = &val
		//p.fulfilmentChannel <- struct{}{}
		close(p.fulfilmentChannel)
		if p.doOnResolved != nil {
			p.doOnResolved(val)
		}
		if p.doOnCompleted != nil {
			p.doOnCompleted()
		}
	})
}

func (p *Progress[T]) reject(err error) {
	p.doneOnce.Do(func() {
		p.outputError = err
		//p.fulfilmentChannel <- struct{}{}
		close(p.fulfilmentChannel)
		if p.doOnRejected != nil {
			p.doOnRejected(err)
		}
		if p.doOnCompleted != nil {
			p.doOnCompleted()
		}
	})
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
