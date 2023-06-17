package promise

import (
	"context"
	"fmt"
	"time"
)

func NewMap[T any](m map[interface{}]interface{}, timeOutLimit time.Duration) FutureMap[T] {
	sizeOfMap := len(m)
	if sizeOfMap == 0 {
		panic(unexpectedCaseError("empty map"))
	}

	fulfilmentChannel := make(chan Output[T], sizeOfMap)

	ctx := context.Background()
	var cancelContext context.Context
	var cancel context.CancelFunc
	if timeOutLimit == 0 {
		cancelContext, cancel = context.WithCancel(ctx)
	} else {
		cancelContext, cancel = context.WithTimeout(ctx, timeOutLimit)
	}
	fm := make(futureMap[T])
	for k, ff := range m {
		switch v := ff.(type) {
		case Future[T]:
			fm[k] = v.setKey(k).setFulfilmentChannel(fulfilmentChannel).setContext(cancelContext, cancel)
		case func() (T, error):
			fm[k] = New(v).TimeOutLimit(timeOutLimit).setKey(k).setFulfilmentChannel(fulfilmentChannel).setContext(cancelContext, cancel)
		case Progress[T]:
			cancel()
			panic(proprietyError("committed promise can not be mapped"))
		default:
			cancel()
			panic(unexpectedCaseError(fmt.Sprintf("unexpected type %T", v)))
		}
	}

	return &fm
}
