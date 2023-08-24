package juck

import (
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
}

func NewScheduledThreadPool(maxGoroutines int) *ScheduledExecutorService {
	p := &ScheduledExecutorService{
		ExecutorService: newThreadPoolExecutor(maxGoroutines),
		scheduledTasks:  make(chan ScheduledTask),
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

func (p *ScheduledExecutorService) ScheduleWithDelay(task func(), delay int64, unit timeunit.TimeUnit) {
	p.scheduledTasks <- ScheduledTask{
		Task:      func() any { task(); return nil },
		ExecuteAt: time.Now().Add(unit.MultiplyWithInt(delay)),
	}
}
