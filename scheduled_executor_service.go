package juck

import (
	"time"

	"github.com/wakamenod/juck/timeunit"
)

type ScheduledExecutorService struct {
	ExecutorService
	cacnelChannel chan any
}

func NewScheduledThreadPool(maxGoroutines int) *ScheduledExecutorService {
	return &ScheduledExecutorService{
		ExecutorService: newThreadPoolExecutor(maxGoroutines),
		cacnelChannel:   make(chan any),
	}
}

func NewSingleScheduledExecutor() *ScheduledExecutorService {
	return NewScheduledThreadPool(1)
}

func (p *ScheduledExecutorService) Shutdown() {
	close(p.cacnelChannel)
	p.ExecutorService.Shutdown()
}

func (p *ScheduledExecutorService) ScheduleWithDelay(task func(), delay int64, unit timeunit.TimeUnit) {
	go func() {
		delay := time.Until(time.Now().Add(unit.MultiplyWithInt(delay)))
		time.Sleep(delay)
		// nolint
		p.SubmitAwait(task)
	}()
}

func (p *ScheduledExecutorService) ScheduleAtFixedRate(task func(), initialDelay, period int64, unit timeunit.TimeUnit) {
	go func() {
		doneInitial := make(chan bool)

		// Initial delay
		<-time.After(unit.MultiplyWithInt(initialDelay))
		go func() {
			f := p.submit(func() any {
				task()
				close(doneInitial)
				return nil
			})
			<-f.result
		}()

		// Setup ticker after initial delay
		ticker := time.NewTicker(unit.MultiplyWithInt(period))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				<-doneInitial
				// nolint
				p.SubmitAwait(task)
			case <-p.cacnelChannel:
				return
			}
		}
	}()
}
