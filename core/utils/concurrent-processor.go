package utils

import (
	"context"
	"sync"
)

type Task struct {
	Id   any
	Data any
}

type ProcessFunc func(ctx context.Context, task Task) Result

type ConcurrentProcessor struct {
	Workers     int
	processFunc ProcessFunc
	TaskChan    chan Task
	ResultChan  chan Result
	wg          sync.WaitGroup
	ctx         context.Context
	cancelFun   context.CancelFunc
}

type Result struct {
	JobID any
	Data  any
	Error error
}

func NewConcurrentProcessor(workers int, processFunc ProcessFunc) *ConcurrentProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConcurrentProcessor{
		Workers:     workers,
		TaskChan:    make(chan Task),
		ResultChan:  make(chan Result),
		ctx:         ctx,
		cancelFun:   cancel,
		processFunc: processFunc,
	}
}

func (cp *ConcurrentProcessor) Start() {
	for range cp.Workers {
		cp.wg.Add(1)
		go cp.worker()
	}
}

func (cp *ConcurrentProcessor) worker() {
	defer cp.wg.Done()

	for {
		select {
		case job, ok := <-cp.TaskChan:
			if !ok {
				return
			}
			result := cp.processFunc(cp.ctx, job)
			cp.ResultChan <- result
		case <-cp.ctx.Done():
			return
		}
	}
}

func (cp *ConcurrentProcessor) Feed(job Task) {
	cp.TaskChan <- job
}

func (cp *ConcurrentProcessor) GetResults() <-chan Result {
	return cp.ResultChan
}

// Stop gracefully shuts down the processor
func (cp *ConcurrentProcessor) Stop() {
	close(cp.TaskChan)
	cp.wg.Wait()
	close(cp.ResultChan)
}

// Cancel immediately stops all workers
func (cp *ConcurrentProcessor) Cancel() {
	cp.cancelFun()
	cp.wg.Wait()
	close(cp.ResultChan)
}
