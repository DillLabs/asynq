package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/examples/basic/tasks"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	hardExitOnInterrupt := flag.Bool("hard-exit-on-interrupt", false, "exit immediately on Ctrl+C without graceful shutdown")
	flag.Parse()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
			},
			// Speed up retry handoff for crash simulation demos.
			RetryDelayFunc: func(n int, err error, t *asynq.Task) time.Duration {
				return 1 * time.Second
			},
			DelayedTaskCheckInterval: 1 * time.Second,
			LeaseDuration:            25 * time.Second,
			HeartbeatInterval:        1 * time.Second,
			RecovererInterval:        1 * time.Second,
			RecovererCutoff:          3 * time.Second,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeWelcomeEmail, tasks.HandleWelcomeEmailTask)
	mux.HandleFunc(tasks.TypeImageResize, tasks.HandleImageResizeTask)

	if *hardExitOnInterrupt {
		if err := srv.Start(mux); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
		log.Printf("worker started in hard-exit mode; press Ctrl+C to exit immediately")

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Printf("received interrupt, exiting immediately without graceful shutdown")
		os.Exit(130)
	}

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
