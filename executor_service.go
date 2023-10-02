package juck

import (
	"sync"
)

type ExecutorService interface {
	Shutdown()
	Submit(task func()) *Future[any]
	SubmitAwait(task func()) (any, error)
	SubmitAny(task func() any) *Future[any]
	SubmitString(task func() string) *Future[string]
	SubmitInt(task func() int) *Future[int]
	SubmitWithAny(task func(), result any) *Future[any]
	SubmitWithString(task func(), result string) *Future[string]
	SubmitWithInt(task func(), result int) *Future[int]
	submit(task func() any) *ft
}

type ThreadPoolExecutor struct {
	taskQueue *BlockingQueue[func()]
	wg        sync.WaitGroup
	closeCh   chan any
}

func newThreadPoolExecutor(maxGoroutines int) ExecutorService {
	p := &ThreadPoolExecutor{
		taskQueue: NewBlockingQueue[func()](),
		closeCh:   make(chan any),
	}

	for i := 0; i < maxGoroutines; i++ {
		go p.worker()
	}

	return p
}

func (p *ThreadPoolExecutor) worker() {
	for {
		select {
		case <-p.closeCh:
			return
		default:
		}
		// TODO let closeCh interrupt the wait?
		// in that case maybe we need take with time e.g TakeWithTimeout
		task := p.taskQueue.Take()
		task()
		p.wg.Done()
	}
}

// func (bq *BlockingQueue[T]) TakeWithTimeout(timeout time.Duration) (T, bool) {
// 	bq.mu.Lock()
// 	defer bq.mu.Unlock()

// 	expireTime := time.Now().Add(timeout)
// 	for len(bq.queue) == 0 {
// 		remainingTime := time.Until(expireTime)
// 		if remainingTime <= 0 {
// 			return zeroValueOf(T), false
// 		}
// 		bq.cond.WaitTimeout(remainingTime)
// 	}

// 	item := bq.queue[0]
// 	bq.queue = bq.queue[1:]
// 	return item, true
// }

func (p *ThreadPoolExecutor) submit(task func() any) *ft {
	resultChan := make(chan any)
	wrappedTask := func() {
		res := task()
		resultChan <- res
		close(resultChan)
	}

	p.wg.Add(1)
	p.taskQueue.Put(wrappedTask)

	return &ft{
		result: resultChan,
	}
}

// func (p *ThreadPoolExecutor) SubmitAsync(task func() any) {
// 	go func() {
// 		f := p.submit(func() any {
// 			task()
// 			return nil
// 		})
// 		<-f.result
// 	}()
// }

func (p *ThreadPoolExecutor) SubmitAwait(task func()) (any, error) {
	f := p.Submit(func() {
		task()
	})
	return f.Await()
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
	close(p.closeCh)
}
