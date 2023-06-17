package promise

import "context"

type cancelableContextProvider interface {
	getCancelContext() context.Context
	getCancelFunction() context.CancelFunc
}
