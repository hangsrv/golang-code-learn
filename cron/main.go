package main

import (
	"log"
	"time"

	"golang-code-learn/cron/custom_scheduled_task"
)

type MyExecutor struct {
}

func (e *MyExecutor) Execute() error {
	return nil
}

func main() {
	taskManager := custom_scheduled_task.NewTaskManager()

	taskManager.AddMultiTask(&MyExecutor{}, 10000, time.Second*1, custom_scheduled_task.WithBefore(func(ctx map[string]interface{}) {
		ctx["start"] = time.Now()
	}), custom_scheduled_task.WithAfter(func(ctx map[string]interface{}) {
		log.Printf("end cost[%+v]", time.Since(ctx["start"].(time.Time)))
	}))

	taskManager.RunAllTasks()
}
