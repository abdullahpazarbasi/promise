package promise

import (
	"fmt"
	"time"
)

type Map interface {
	AssureCommitments() Map
	Cancel()
	Await() *map[interface{}]struct {
		out interface{}
		err error
	}
	Race() (key interface{}, out interface{}, err error)
}

type promiseMap map[interface{}]Promise[interface{}]

func NewMap(rm map[interface{}]interface{}, defaultTimeOutLimit time.Duration) Map {
	if len(rm) == 0 {
		panic(unexpectedCaseError("empty map"))
	}
	pm := make(promiseMap)
	for k, progressOrFunction := range rm {
		switch v := progressOrFunction.(type) {
		case Promise[any]:
			pm[k] = v.TimeOutLimit(defaultTimeOutLimit)
		case func() (interface{}, error):
			pm[k] = New(v).TimeOutLimit(defaultTimeOutLimit)
		default:
			panic(unexpectedCaseError(fmt.Sprintf("unexpected type %T", v)))
		}
	}

	return &pm
}

func (m *promiseMap) AssureCommitments() Map {
	nm := make(promiseMap)
	for k, p := range *m {
		nm[k] = p.Commit()
	}

	return &nm
}

func (m *promiseMap) Cancel() {
	for _, p := range *m {
		switch v := p.(type) {
		case *Progress[any]:
			v.Cancel()
		}
	}
}

func (m *promiseMap) Await() *map[interface{}]struct {
	out interface{}
	err error
} {
	om := map[interface{}]struct {
		out interface{}
		err error
	}{}
	for k, p := range *m {
		select {
		case <-p.getContext().Done():
			o, e := p.abandon(p.getContext().Err())
			om[k] = struct {
				out interface{}
				err error
			}{
				out: o,
				err: e,
			}
		case r := <-p.getFulfilmentChannel():
			o, e := p.fulfil(r)
			om[k] = struct {
				out interface{}
				err error
			}{
				out: o,
				err: e,
			}
		}
	}

	return &om
}

func (m *promiseMap) Race() (key interface{}, out interface{}, err error) {
	//TODO implement me
	panic("implement me")
}
