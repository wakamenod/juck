package juck

import "errors"

// inner future to get result as any type
type ft struct {
	result chan any
}

type Future[T any] struct {
	*ft
}

func (f *Future[T]) Await() (T, error) {
	if value, ok := (<-f.result).(T); ok {
		return value, nil
	} else {
		return *new(T), errors.New("the result is not of the expected type")
	}
}

func (f *Future[T]) Ft() chan any {
	return f.result
}
