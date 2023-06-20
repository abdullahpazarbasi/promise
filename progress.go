package promise

import (
	"context"
	"fmt"
	"sync"
)

// Progress is promise which is committed and can be awaited or canceled
type Progress[T any] interface {
	Cancel()
	Await() (T, error)
	await() (T, error)
	getFulfilmentChannel() chan Output[T]
	abandon()
	leave()
	eliminate()
}

type progress[T any] struct {
	doOnResolved      func(T)
	doOnRejected      func(error)
	doOnCanceled      func()
	doOnTimedOut      func()
	doFinally         func(event)
	fulfilmentChannel chan Output[T]
	cancelableContext context.Context
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
	var z T
	select {
	case <-p.cancelableContext.Done():
		e := p.cancelableContext.Err()
		switch e {
		case context.Canceled:
			p.abandon()
		case context.DeadlineExceeded:
			p.leave()
		}

		return z, e
	case r, o := <-p.fulfilmentChannel:
		defer p.cancel()
		if o {
			defer close(p.fulfilmentChannel)
		}

		if r == nil {
			return z, nil
		}

		return r.Payload(), r.Error()
	}
}

func (p *progress[T]) getFulfilmentChannel() chan Output[T] {
	return p.fulfilmentChannel
}

func (p *progress[T]) getContext() context.Context {
	return p.cancelableContext
}

func (p *progress[T]) getCancelFunction() context.CancelFunc {
	return p.cancel
}

func (p *progress[T]) abandon() {
	p.doneOnce.Do(func() {
		if p.doOnCanceled != nil {
			p.doOnCanceled()
		}
		if p.doFinally != nil {
			p.doFinally(EventCanceled)
		}
	})
}

func (p *progress[T]) leave() {
	p.doneOnce.Do(func() {
		if p.doOnTimedOut != nil {
			p.doOnTimedOut()
		}
		if p.doFinally != nil {
			p.doFinally(EventTimedOut)
		}
	})
}

func (p *progress[T]) eliminate() {
	p.doneOnce.Do(func() {
		if p.doFinally != nil {
			p.doFinally(EventEliminated)
		}
	})
}

func (p *progress[T]) handleProbablePanic() {
	e := recover()
	if e == nil {
		return
	}

	switch v := e.(type) {
	case error:
		p.reject(v)
	default:
		p.reject(fmt.Errorf("%+v", v))
	}
}

func (p *progress[T]) resolve(pay T) {
	p.doneOnce.Do(func() {
		p.fulfilmentChannel <- newOutput[T](pay, nil).setKey(p.key)
		if p.doOnResolved != nil {
			p.doOnResolved(pay)
		}
		if p.doFinally != nil {
			p.doFinally(EventResolved)
		}
	})
}

func (p *progress[T]) reject(err error) {
	p.doneOnce.Do(func() {
		var z T
		p.fulfilmentChannel <- newOutput[T](z, err).setKey(p.key)
		if p.doOnRejected != nil {
			p.doOnRejected(err)
		}
		if p.doFinally != nil {
			p.doFinally(EventRejected)
		}
	})
}
