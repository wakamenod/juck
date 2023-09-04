package juck

import (
	"fmt"
	"sync"
)

type ExecutorService interface {
	Shutdown()
	Submit(task func()) *Future[any]
	SubmitAny(task func() any) *Future[any]
	SubmitString(task func() string) *Future[string]
	SubmitInt(task func() int) *Future[int]
	SubmitWithAny(task func(), result any) *Future[any]
	SubmitWithString(task func(), result string) *Future[string]
	SubmitWithInt(task func(), result int) *Future[int]
	submit(task func() any) *ft
}

type ThreadPoolExecutor struct {
	taskQueue    chan func()
	workerNumber int
	wg           sync.WaitGroup
}

func newThreadPoolExecutor(maxGoroutines int) ExecutorService {
	p := &ThreadPoolExecutor{
		taskQueue:    make(chan func()),
		workerNumber: maxGoroutines,
	}

	for i := 0; i < maxGoroutines; i++ {
		go p.worker()
	}

	return p
}

func (p *ThreadPoolExecutor) worker() {
	for task := range p.tasks {
		task()
		p.wg.Done()
	}
}

func (p *ThreadPoolExecutor) submit(task func() any) *ft {
	fmt.Println("@@@@@@@@@@@ submit")

	resultChan := make(chan any)
	wrappedTask := func() {
		res := task()
		fmt.Println("11111111111111111")
		resultChan <- res
		fmt.Println("2222222222222222222")
		close(resultChan)
	}

	fmt.Println("@@@@@@@@@@@ wrappedTask start")
	p.tasks <- wrappedTask
	fmt.Println("@@@@@@@@@@@ wrappedTask end")
	p.wg.Add(1)

	return &ft{
		result: resultChan,
	}
}

func NewFixedThreadPool(nThreads int) ExecutorService {
	return newThreadPoolExecutor(nThreads)
}

func NewSingleThreadPool() ExecutorService {
	return newThreadPoolExecutor(1)
}

func (p *ThreadPoolExecutor) SubmitAny(task func() any) *Future[any] {
	f := p.submit(func() any { return task() })
	return &Future[any]{f}
}

func (p *ThreadPoolExecutor) SubmitString(task func() string) *Future[string] {
	f := p.submit(func() any { return task() })
	return &Future[string]{f}
}

func (p *ThreadPoolExecutor) SubmitInt(task func() int) *Future[int] {
	f := p.submit(func() any { return task() })
	return &Future[int]{f}
}

func (p *ThreadPoolExecutor) Submit(task func()) *Future[any] {
	f := p.submit(func() any { task(); return nil })
	return &Future[any]{f}
}

func (p *ThreadPoolExecutor) SubmitWithString(task func(), result string) *Future[string] {
	f := p.submit(func() any { task(); return result })
	return &Future[string]{f}
}

func (p *ThreadPoolExecutor) SubmitWithInt(task func(), result int) *Future[int] {
	f := p.submit(func() any { task(); return result })
	return &Future[int]{f}
}

func (p *ThreadPoolExecutor) SubmitWithAny(task func(), result any) *Future[any] {
	f := p.submit(func() any { task(); return result })
	return &Future[any]{f}
}

func (p *ThreadPoolExecutor) Shutdown() {
	close(p.tasks)
}
