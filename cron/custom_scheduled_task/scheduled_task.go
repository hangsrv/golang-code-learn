package custom_scheduled_task

import (
	"context"
	"log"
	"time"
)

type Executor interface {
	Execute() error
}

type Task struct {
	executor Executor

	ctx map[string]interface{}

	interval time.Duration
	once     bool

	before  []func(ctx map[string]interface{})
	after   []func(ctx map[string]interface{})
	recover func(interface{})
}

func newTask(executor Executor, interval time.Duration, options ...func(*Task)) *Task {
	task := &Task{
		executor: executor,
		interval: interval,
		ctx:      make(map[string]interface{}),
		recover: func(r interface{}) {
			log.Printf("Panic recovered: %v", r)
		},
	}

	for _, option := range options {
		option(task)
	}

	return task
}

func (t *Task) run(ctx context.Context) {
	if t.once {
		t.runOnce(ctx)
	} else {
		t.runPeriodically(ctx)
	}
}

func (t *Task) runPeriodically(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping task...")
			return
		default:
			t.executeCore()
		}
	}
}

func (t *Task) runOnce(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Println("Stopping task...")
		return
	default:
		t.executeCore()
	}
}
func (t *Task) executeCore() {
	time.Sleep(t.interval)
	t.execute()
}
func (t *Task) execute() {
	defer func() {
		if r := recover(); r != nil {
			t.recover(r)
		}
	}()

	for _, beforeFunc := range t.before {
		beforeFunc(t.ctx)
	}

	if err := t.executor.Execute(); err != nil {
		log.Printf("Task execution failed: %v", err)
	}

	for _, afterFunc := range t.after {
		afterFunc(t.ctx)
	}
}

func WithBefore(beforeFuncs ...func(ctx map[string]interface{})) func(*Task) {
	return func(task *Task) {
		task.before = append(task.before, beforeFuncs...)
	}
}

func WithAfter(afterFuncs ...func(ctx map[string]interface{})) func(*Task) {
	return func(task *Task) {
		task.after = append(task.after, afterFuncs...)
	}
}

func WithRecover(recoverFunc func(interface{})) func(*Task) {
	return func(task *Task) {
		task.recover = recoverFunc
	}
}

func WithOnce() func(*Task) {
	return func(task *Task) {
		task.once = true
	}
}
