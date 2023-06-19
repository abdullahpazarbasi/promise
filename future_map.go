package promise

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type FutureMap[T any] interface {
	TimeOutLimit(timeOutLimit time.Duration) FutureMap[T]
	Context(ctx context.Context) FutureMap[T]
	Commit() ProgressMap[T]
	Await() *map[interface{}]Output[T]
	Race() (key interface{}, pay T, err error)
	commit() ProgressMap[T]
}

type futureMap[T any] struct {
	sourceMap           *map[interface{}]interface{}
	timeOutLimit        time.Duration
	context             context.Context
	timeOutLimitSetOnce sync.Once
	contextSetOnce      sync.Once
	committedOnce       sync.Once
}

func (fm *futureMap[T]) TimeOutLimit(timeOutLimit time.Duration) FutureMap[T] {
	fm.timeOutLimitSetOnce.Do(func() {
		fm.timeOutLimit = timeOutLimit
	})

	return fm
}

func (fm *futureMap[T]) Context(ctx context.Context) FutureMap[T] {
	fm.contextSetOnce.Do(func() {
		fm.context = ctx
	})

	return fm
}

func (fm *futureMap[T]) Commit() ProgressMap[T] {
	return fm.commit()
}

func (fm *futureMap[T]) Await() *map[interface{}]Output[T] {
	return fm.commit().await()
}

func (fm *futureMap[T]) Race() (key interface{}, pay T, err error) {
	return fm.commit().race()
}

func (fm *futureMap[T]) commit() ProgressMap[T] {
	var pm progressMap[T]
	fm.committedOnce.Do(func() {
		pm.data = map[interface{}]Progress[T]{}
		sizeOfMap := len(*fm.sourceMap)
		pm.fulfilmentChannel = make(chan Output[T], sizeOfMap)
		fm.contextSetOnce.Do(func() {
			fm.context = context.Background()
		})
		if fm.timeOutLimit == 0 {
			pm.cancelableContext, pm.cancel = context.WithCancel(fm.context)
		} else {
			pm.cancelableContext, pm.cancel = context.WithTimeout(fm.context, fm.timeOutLimit)
		}
		data := map[interface{}]Future[T]{}
		var cancelableContext context.Context
		var cancel context.CancelFunc
		for k, ff := range *fm.sourceMap {
			switch v := ff.(type) {
			case Future[T], func(context.Context) (T, error):
				var f Future[T]
				switch vv := v.(type) {
				case Future[T]:
					f = vv
				case func(context.Context) (T, error):
					f = New(vv).TimeOutLimit(fm.timeOutLimit)
				}
				timeOutLimit := f.getTimeOutLimit()
				if timeOutLimit == 0 {
					cancelableContext, cancel = context.WithCancel(pm.cancelableContext)
				} else {
					cancelableContext, cancel = context.WithTimeout(pm.cancelableContext, timeOutLimit)
				}
				data[k] = f.setKey(k).setFulfilmentChannel(pm.fulfilmentChannel).setCancelableContextAndCancelFunction(cancelableContext, cancel)
			case Progress[T]:
				panic(proprietyError("committed promise can not be mapped"))
			default:
				panic(unexpectedCaseError(fmt.Sprintf("unexpected type %T", v)))
			}
		}
		for k, p := range data {
			pm.data[k] = p.commit()
		}
	})

	return &pm
}
