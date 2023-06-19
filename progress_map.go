package promise

import "context"

type ProgressMap[T any] interface {
	Cancel()
	Await() *map[interface{}]Output[T]
	Race() (key interface{}, pay T, err error)
	await() *map[interface{}]Output[T]
	race() (key interface{}, pay T, err error)
}

type progressMap[T any] struct {
	data              map[interface{}]Progress[T]
	fulfilmentChannel chan Output[T]
	cancelableContext context.Context
	cancel            context.CancelFunc
}

func (pm *progressMap[T]) Cancel() {
	defer pm.cancel()
}

func (pm *progressMap[T]) Await() *map[interface{}]Output[T] {
	return pm.await()
}

func (pm *progressMap[T]) Race() (key interface{}, pay T, err error) {
	return pm.race()
}

func (pm *progressMap[T]) await() *map[interface{}]Output[T] {
	defer pm.cancel()

	var z T

	om := map[interface{}]Output[T]{}
	var d bool
	var o bool
	var nio Output[T]
	som := len(pm.data)
	for {
		for k, p := range pm.data {
			c := p.(cancelableContextProvider).getContext()
			select {
			case <-c.Done():
				_, ok := om[k]
				if !ok {
					e := c.Err()
					switch e {
					case context.Canceled:
						p.abandon()
					case context.DeadlineExceeded:
						p.leave()
					}
					om[k] = newOutput(z, e)
				}
				if len(om) >= som {
					return &om
				}
			case nio, o = <-pm.fulfilmentChannel:
				if o {
					if !d {
						d = true
						defer close(pm.fulfilmentChannel)
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
}

func (pm *progressMap[T]) race() (key interface{}, pay T, err error) {
	defer pm.cancel()

	om := map[interface{}]bool{}
	var d bool
	som := len(pm.data)
	for {
		for k, p := range pm.data {
			c := p.(cancelableContextProvider).getContext()
			select {
			case <-c.Done():
				key = nil
				err = c.Err()
				switch err {
				case context.Canceled:
					p.abandon()
				case context.DeadlineExceeded:
					p.leave()
				}
				om[k] = true
				if len(om) >= som {
					return
				}
			case nio, o := <-pm.fulfilmentChannel:
				if o {
					if !d {
						d = true
						defer close(pm.fulfilmentChannel)
					}
				}
				key = nio.(keyProvider).Key()
				pay = nio.Payload()
				err = nio.Error()
				om[key] = true
				for kk, pp := range pm.data {
					if kk != key {
						pp.eliminate()
						om[kk] = true
					}
				}
				if len(om) >= som {
					return
				}
			}
		}
	}
}
