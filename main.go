package main

import (
	"fmt"
	"sync"
)

type Runnable interface {
	Run()
}

type Executor interface {
	Execute(r Runnable)
}

type Pool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

// inner future to get result as any type
type ft struct {
	result chan any
}

type Future[T any] struct {
	*ft
}

// Get is to get something
func (f *Future[T]) Get() T {
	return (<-f.result).(T)
}

func NewPool(maxGoroutines int) *Pool {
	p := &Pool{
		tasks: make(chan func()),
	}

	for i := 0; i < maxGoroutines; i++ {
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	for task := range p.tasks {
		task()
		p.wg.Done()
	}
}

func (p *Pool) submit(task func() any) *ft {
	resultChan := make(chan any)
	wrappedTask := func() {
		res := task()
		resultChan <- res
		close(resultChan)
	}

	p.tasks <- wrappedTask
	p.wg.Add(1)

	return &ft{
		result: resultChan,
	}
}

func (p *Pool) SubmitAny(task func() any) *Future[any] {
	f := p.submit(func() any { return task() })
	return &Future[any]{f}
}

func (p *Pool) SubmitString(task func() string) *Future[string] {
	f := p.submit(func() any { return task() })
	return &Future[string]{f}
}

func (p *Pool) SubmitInt(task func() int) *Future[int] {
	f := p.submit(func() any { return task() })
	return &Future[int]{f}
}

func (p *Pool) Submit(task func()) *Future[any] {
	f := p.submit(func() any { task(); return nil })
	return &Future[any]{f}
}

func (p *Pool) SubmitWithResultString(task func(), result string) *Future[string] {
	f := p.submit(func() any { task(); return result })
	return &Future[string]{f}
}

func (p *Pool) SubmitWithResultInt(task func(), result int) *Future[int] {
	f := p.submit(func() any { task(); return result })
	return &Future[int]{f}
}

func (p *Pool) Shutdown() {
	close(p.tasks)
}

func main() {
	pool := NewPool(3)

	stringFuture := pool.SubmitString(func() string {
		return "Hello, world!"
	})

	intFuture := pool.SubmitInt(func() int {
		return 42
	})

	pool.Shutdown()

	fmt.Println(stringFuture.Get())
	fmt.Println(intFuture.Get())
}
