package main

import (
	"log"
	"time"

	"github.com/DillLabs/asynq"
	"github.com/DillLabs/asynq/examples/scheduler/tasks"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("could not load location: %v", err)
	}

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{
			Location: loc,
			PreEnqueueFunc: func(task *asynq.Task, opts []asynq.Option) {
				log.Printf("about to enqueue task: type=%s", task.Type())
			},
			PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
				if err != nil {
					log.Printf("enqueue error: %v", err)
					return
				}
				log.Printf("enqueued periodic task: id=%s queue=%s", info.ID, info.Queue)
			},
		},
	)

	task, err := tasks.NewReportTask(time.Now().Format("2006-01-02"))
	if err != nil {
		log.Fatalf("could not create report task: %v", err)
	}

	entryID, err := scheduler.Register(
		//"* * * * *",
		"@every 6s",
		task,
		asynq.Queue("default"),
		asynq.MaxRetry(3),
	)
	if err != nil {
		log.Fatalf("could not register scheduler entry: %v", err)
	}
	log.Printf("registered scheduler entry: id=%s", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
	}
}
