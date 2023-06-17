package promise

import "context"

type ProgressMap[T any] interface {
	Cancel()
	Await() *map[interface{}]Output[T]
	Race() (key interface{}, pay T, err error)
	await() *map[interface{}]Output[T]
}

type progressMap[T any] map[interface{}]Progress[T]

func (pm *progressMap[T]) Cancel() {
	var cancel context.CancelFunc
	var fc chan Output[T]
	for _, p := range *pm {
		cancel = p.(cancelableContextProvider).getCancelFunction()
		fc = p.getFulfilmentChannel()
		break
	}
	defer cancel()
	defer close(fc)
}

func (pm *progressMap[T]) Await() *map[interface{}]Output[T] {
	return pm.await()
}

func (pm *progressMap[T]) await() *map[interface{}]Output[T] {
	var ky interface{}
	var ps Progress[T]
	var cx context.Context
	var cancel context.CancelFunc
	var fc chan Output[T]
	for ky, ps = range *pm {
		cx = ps.(cancelableContextProvider).getCancelContext()
		cancel = ps.(cancelableContextProvider).getCancelFunction()
		fc = ps.getFulfilmentChannel()
		break
	}

	defer close(fc)
	defer cancel()

	om := map[interface{}]Output[T]{}
	var nextInternalOutput Output[T]
	sizeOfMap := len(*pm)
	for {
		select {
		case <-cx.Done():
			for ky, ps = range *pm {
				om[ky] = newOutput(ps.abandon())
			}

			return &om
		case nextInternalOutput = <-fc:
			om[nextInternalOutput.(keyProvider).Key()] = newOutput(
				nextInternalOutput.Payload(),
				nextInternalOutput.Error(),
			)
			if len(om) >= sizeOfMap {
				return &om
			}
		}
	}
}

func (pm *progressMap[T]) Race() (key interface{}, pay T, err error) {
	// todo:
	panic("not implemented yet")
}
