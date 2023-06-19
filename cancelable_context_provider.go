package promise

import "context"

type cancelableContextProvider interface {
	getContext() context.Context
	getCancelFunction() context.CancelFunc
}
