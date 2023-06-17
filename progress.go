package promise

import (
	"context"
	"fmt"
	"sync"
)

type Progress[T any] interface {
	Cancel()
	Await() (T, error)
	await() (T, error)
	getFulfilmentChannel() chan Output[T]
	abandon() (T, error)
}

type progress[T any] struct {
	doOnResolved      func(T)
	doOnRejected      func(error)
	doOnCanceled      func()
	doOnTimedOut      func()
	doFinally         func(event)
	fulfilmentChannel chan Output[T]
	context           context.Context
	cancel            context.CancelFunc
	key               interface{}
	doneOnce          sync.Once
}

func (p *progress[T]) Cancel() {
	defer p.cancel()
}

func (p *progress[T]) Await() (T, error) {
	return p.await()
}

func (p *progress[T]) await() (T, error) {
	select {
	case <-p.context.Done():
		return p.abandon()
	case r, open := <-p.fulfilmentChannel:
		defer p.cancel()
		if open {
			defer close(p.fulfilmentChannel)
		}

		if r == nil {
			var zeroT T

			return zeroT, nil
		}

		return r.Payload(), r.Error()
	}
}

func (p *progress[T]) getFulfilmentChannel() chan Output[T] {
	return p.fulfilmentChannel
}

func (p *progress[T]) getCancelContext() context.Context {
	return p.context
}

func (p *progress[T]) getCancelFunction() context.CancelFunc {
	return p.cancel
}

func (p *progress[T]) abandon() (T, error) {
	var zeroT T
	err := p.context.Err()

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

	return zeroT, err
}

func (p *progress[T]) handleProbablePanic() {
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

func (p *progress[T]) resolve(out T) {
	p.doneOnce.Do(func() {
		p.fulfilmentChannel <- newOutput[T](out, nil).setKey(p.key)
		if p.doOnResolved != nil {
			p.doOnResolved(out)
		}
		if p.doFinally != nil {
			p.doFinally(EventResolved)
		}
	})
}

func (p *progress[T]) reject(err error) {
	p.doneOnce.Do(func() {
		var zeroT T
		p.fulfilmentChannel <- newOutput[T](zeroT, err).setKey(p.key)
		if p.doOnRejected != nil {
			p.doOnRejected(err)
		}
		if p.doFinally != nil {
			p.doFinally(EventRejected)
		}
	})
}
