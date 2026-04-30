package main

import (
	"flag"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/examples/basic/tasks"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	taskType := flag.String("task", "all", "task type to enqueue: all|email|delayed-email|image")
	flag.Parse()

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	if *taskType != "all" && *taskType != "email" && *taskType != "delayed-email" && *taskType != "image" {
		log.Fatalf("invalid -task value %q, must be one of: all|email|delayed-email|image", *taskType)
	}

	if *taskType == "all" || *taskType == "email" {
		emailTask, err := tasks.NewWelcomeEmailTask(42, "user@example.com")
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}

		info, err := client.Enqueue(
			emailTask,
			asynq.Queue("critical"),
			asynq.MaxRetry(5),
			asynq.Timeout(20*time.Second),
		)
		if err != nil {
			log.Fatalf("could not enqueue immediate task: %v", err)
		}
		log.Printf("enqueued immediate task: id=%s queue=%s", info.ID, info.Queue)
	}

	if *taskType == "all" || *taskType == "delayed-email" {
		delayedEmailTask, err := tasks.NewWelcomeEmailTask(7, "delayed@example.com")
		if err != nil {
			log.Fatalf("could not create delayed task: %v", err)
		}
		info, err := client.Enqueue(
			delayedEmailTask,
			asynq.Queue("default"),
			asynq.MaxRetry(5),
			asynq.Timeout(20*time.Second),
			asynq.ProcessIn(10*time.Second),
		)
		if err != nil {
			log.Fatalf("could not enqueue delayed task: %v", err)
		}
		log.Printf("enqueued delayed task: id=%s queue=%s", info.ID, info.Queue)
	}

	if *taskType == "all" || *taskType == "image" {
		imageTask, err := tasks.NewImageResizeTask("https://example.com/assets/image.jpg")
		if err != nil {
			log.Fatalf("could not create image task: %v", err)
		}
		info, err := client.Enqueue(
			imageTask,
			asynq.Queue("default"),
			asynq.MaxRetry(3),
			asynq.Timeout(10*time.Second),
		)
		if err != nil {
			log.Fatalf("could not enqueue image task: %v", err)
		}
		log.Printf("enqueued timeout-demo task: id=%s queue=%s", info.ID, info.Queue)
	}
}
