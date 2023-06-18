package promise

import "context"

type ProgressMap[T any] interface {
	Cancel()
	Await() *map[interface{}]Output[T]
	Race() (key interface{}, pay T, err error)
	await() *map[interface{}]Output[T]
	race() (key interface{}, pay T, err error)
}

type progressMap[T any] map[interface{}]Progress[T]

func (pm *progressMap[T]) Cancel() {
	var cancel context.CancelFunc
	for _, p := range *pm {
		cancel = p.(cancelableContextProvider).getCancelFunction()
		break
	}
	defer cancel()
}

func (pm *progressMap[T]) Await() *map[interface{}]Output[T] {
	return pm.await()
}

func (pm *progressMap[T]) Race() (key interface{}, pay T, err error) {
	return pm.race()
}

func (pm *progressMap[T]) await() *map[interface{}]Output[T] {
	var k interface{}
	var p Progress[T]
	var z T
	var c context.Context
	var cancel context.CancelFunc
	var fc chan Output[T]
	for k, p = range *pm {
		c = p.(cancelableContextProvider).getCancelContext()
		cancel = p.(cancelableContextProvider).getCancelFunction()
		fc = p.getFulfilmentChannel()
		break
	}

	defer cancel()

	om := map[interface{}]Output[T]{}
	var e error
	var d bool
	var o bool
	var nio Output[T]
	som := len(*pm)
	for {
		select {
		case <-c.Done():
			switch c.Err() {
			case context.Canceled:
				e = canceledError("manually canceled")
			case context.DeadlineExceeded:
				e = timedOutError("timed-out")
			}
			for k, p = range *pm {
				_, ok := om[k]
				if !ok {
					switch e.(type) {
					case canceledError:
						p.abandon(e)
					case timedOutError:
						p.leave(e)
					}
					om[k] = newOutput(z, e)
				}
			}

			return &om
		case nio, o = <-fc:
			if o {
				if !d {
					d = true
					defer close(fc)
				}
			}
			om[nio.(keyProvider).Key()] = newOutput(
				nio.Payload(),
				nio.Error(),
			)
			if len(om) >= som {
				return &om
			}
		}
	}
}

func (pm *progressMap[T]) race() (key interface{}, pay T, err error) {
	var k interface{}
	var p Progress[T]
	var c context.Context
	var cancel context.CancelFunc
	var fc chan Output[T]
	for _, p = range *pm {
		c = p.(cancelableContextProvider).getCancelContext()
		cancel = p.(cancelableContextProvider).getCancelFunction()
		fc = p.getFulfilmentChannel()
		break
	}

	defer cancel()

	var d bool
	var o bool
	var nio Output[T]
	for {
		select {
		case <-c.Done():
			switch c.Err() {
			case context.Canceled:
				err = canceledError("manually canceled")
			case context.DeadlineExceeded:
				err = timedOutError("timed-out")
			}
			key = nil
			for _, p = range *pm {
				switch err.(type) {
				case canceledError:
					p.abandon(err)
				case timedOutError:
					p.leave(err)
				}
			}
			return
		case nio, o = <-fc:
			if o {
				if !d {
					d = true
					defer close(fc)
				}
			}
			key = nio.(keyProvider).Key()
			pay = nio.Payload()
			err = nio.Error()
			for k, p = range *pm {
				if k != key {
					p.eliminate()
				}
			}
			return
		}
	}
}
