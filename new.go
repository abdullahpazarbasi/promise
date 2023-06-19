package promise

import (
	"context"
	"sync"
)

func New[T any](async func(ctx context.Context) (T, error)) Future[T] {
	return &future[T]{
		async:                    async,
		timeOutLimit:             0,
		doOnResolvedSetOnce:      sync.Once{},
		doOnRejectedSetOnce:      sync.Once{},
		doOnCanceledSetOnce:      sync.Once{},
		doOnTimedOutSetOnce:      sync.Once{},
		doFinallySetOnce:         sync.Once{},
		fulfilmentChannelSetOnce: sync.Once{},
		cancelableContextSetOnce: sync.Once{},
		committedOnce:            sync.Once{},
	}
}
