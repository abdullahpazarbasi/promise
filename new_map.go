package promise

import (
	"sync"
)

func NewMap[T any](data map[interface{}]interface{}) FutureMap[T] {
	if len(data) == 0 {
		panic(unexpectedCaseError("empty map"))
	}

	return &futureMap[T]{
		sourceMap:           &data,
		timeOutLimit:        0,
		timeOutLimitSetOnce: sync.Once{},
		contextSetOnce:      sync.Once{},
		committedOnce:       sync.Once{},
	}
}
