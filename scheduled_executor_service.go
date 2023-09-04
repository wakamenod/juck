package juck

import (
	"fmt"
	"time"

	"github.com/wakamenod/juck/timeunit"
)

type ScheduledTask struct {
	Task      func() any
	ExecuteAt time.Time
}

type ScheduledExecutorService struct {
	ExecutorService
	scheduledTasks chan ScheduledTask
	cacnelChannel  chan any
}

func NewScheduledThreadPool(maxGoroutines int) *ScheduledExecutorService {
	p := &ScheduledExecutorService{
		ExecutorService: newThreadPoolExecutor(maxGoroutines),
		scheduledTasks:  make(chan ScheduledTask),
		cacnelChannel:   make(chan any),
	}

	go p.scheduleWorker()
	return p
}

func NewSingleScheduledExecutor() ExecutorService {
	return NewScheduledThreadPool(1)
}

func (p *ScheduledExecutorService) scheduleWorker() {
	for st := range p.scheduledTasks {
		go func(st ScheduledTask) {
			delay := time.Until(st.ExecuteAt)
			time.Sleep(delay)
			p.submit(st.Task)
		}(st)
	}
}

func (p *ScheduledExecutorService) Shutdown() {
	close(p.cacnelChannel)
	p.ExecutorService.Shutdown()
}

func (p *ScheduledExecutorService) ScheduleWithDelay(task func(), delay int64, unit timeunit.TimeUnit) {
	p.scheduledTasks <- ScheduledTask{
		Task:      func() any { task(); return nil },
		ExecuteAt: time.Now().Add(unit.MultiplyWithInt(delay)),
	}
}

func (p *ScheduledExecutorService) ScheduleAtFixedRate(task func(), initialDelay, period int64, unit timeunit.TimeUnit) {
	go func() {
		doneInitial := make(chan bool)

		// Initial delay
		<-time.After(unit.MultiplyWithInt(initialDelay))
		p.submit(func() any {
			task()
			close(doneInitial)
			return nil
		})

		// Setup ticker after initial delay
		ticker := time.NewTicker(unit.MultiplyWithInt(period))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("111111111111111111111")
				<-doneInitial
				fmt.Println("222222222222222222")
				p.submitWait(task)
				fmt.Println("333333333333333333333")
			case <-p.cacnelChannel:
				return
			}
		}
	}()
}

// submitWait blocks until submitted task is finished
func (p *ScheduledExecutorService) submitWait(task func()) {
	fmt.Println("submitWait start")
	done := make(chan bool)

	p.submit(func() any {
		fmt.Println("about to call task()")
		task()
		done <- true
		return nil
	})

	<-done
	fmt.Println("submitWait end")
}
