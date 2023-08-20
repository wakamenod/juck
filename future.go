package juck

// inner future to get result as any type
type ft struct {
	result chan any
}

type Future[T any] struct {
	*ft
}

func (f *Future[T]) Get() T {
	return (<-f.result).(T)
}

func (f *Future[T]) Ft() chan any {
	return f.result
}
