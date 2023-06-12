package promise

import (
	"fmt"
	"time"
)

type Map interface {
	AssureCommitments() Map
	Cancel()
	Await()
	Race()
}

type promiseMap map[interface{}]Promise[interface{}]

func NewMap(rm map[interface{}]interface{}, defaultTimeOutLimit time.Duration) Map {
	if len(rm) == 0 {
		panic(unexpectedCaseError("empty map"))
	}
	pm := make(promiseMap)
	for k, progressOrFunction := range rm {
		switch v := progressOrFunction.(type) {
		case *Future[interface{}]:
			pm[k] = v.TimeOutLimit(defaultTimeOutLimit)
		case *Progress[interface{}]:
			pm[k] = v
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
		switch v := p.(type) {
		case *Future[any]:
			nm[k] = v.Commit()
		default:
			nm[k] = v
		}
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

func (m *promiseMap) Await() {
	//TODO implement me
	panic("implement me")
}

func (m *promiseMap) Race() {
	//TODO implement me
	panic("implement me")
}
