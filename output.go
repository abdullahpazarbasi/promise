package promise

import "reflect"

// Output is interface of structure for asynchronous function 's return values
type Output[T any] interface {
	Payload() T
	Error() error
	setKey(key interface{}) Output[T]
}

type output[T any] struct {
	payload *T
	error   error
	key     interface{}
}

func newOutput[T any](pay T, err error) Output[T] {
	var zeroT T

	if reflect.DeepEqual(pay, zeroT) {
		return &output[T]{
			payload: nil,
			error:   err,
		}
	}

	return &output[T]{
		payload: &pay,
		error:   err,
	}
}

func (o *output[T]) Payload() T {
	if o.payload == nil {
		var zeroT T

		return zeroT
	}

	return *o.payload
}

func (o *output[T]) Error() error {
	return o.error
}

func (o *output[T]) Key() interface{} {
	return o.key
}

func (o *output[T]) setKey(key interface{}) Output[T] {
	o.key = key

	return o
}
