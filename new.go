package promise

import (
	"sync"
)

func New[T any](async func() (T, error)) Future[T] {
	return &future[T]{
		async:                    async,
		timeOutLimit:             0,
		doOnResolvedSetOnce:      sync.Once{},
		doOnRejectedSetOnce:      sync.Once{},
		doOnCanceledSetOnce:      sync.Once{},
		doOnTimedOutSetOnce:      sync.Once{},
		doFinallySetOnce:         sync.Once{},
		fulfilmentChannelSetOnce: sync.Once{},
		cancelContextSetOnce:     sync.Once{},
		committedOnce:            sync.Once{},
	}
}
