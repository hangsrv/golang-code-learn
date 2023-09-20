package custom_scheduled_task

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TaskManager struct {
	tasks  []*Task
	ctx    context.Context
	cancel context.CancelFunc
	mutex  *sync.Mutex
	wg     *sync.WaitGroup
}

func NewTaskManager() *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	tm := &TaskManager{
		ctx:    ctx,
		cancel: cancel,
		tasks:  make([]*Task, 0),
		mutex:  &sync.Mutex{},
		wg:     &sync.WaitGroup{},
	}
	return tm
}

func (tm *TaskManager) AddTask(executor Executor, interval time.Duration, options ...func(*Task)) {
	task := newTask(executor, interval, options...)
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.tasks = append(tm.tasks, task)
}

func (tm *TaskManager) AddMultiTask(executor Executor, concurrence int, interval time.Duration, options ...func(*Task)) {
	for i := 0; i < concurrence; i++ {
		tm.AddTask(executor, interval, options...)
	}
}

func (tm *TaskManager) RunAllTasks() {
	for _, task := range tm.tasks {
		tm.wg.Add(1)
		go func(t *Task) {
			defer tm.wg.Done()
			t.run(tm.ctx)
		}(task)
	}
	log.Println("All tasks started")
	tm.listenSignals()
}

func (tm *TaskManager) stop() {
	tm.cancel()
	tm.wg.Wait()
	log.Println("All tasks stoped")
}

func (tm *TaskManager) listenSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	tm.stop()
}
