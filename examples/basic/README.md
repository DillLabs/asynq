# Basic Asynq Demo

This demo shows a producer/worker setup using the `github.com/DillLabs/asynq` package in this repository, including delayed tasks, queue priority, and timeout/retry behavior.

## Prerequisites

- Go installed
- Redis running on `127.0.0.1:6379`

## Files

- `tasks/tasks.go`: task type, payload, and handlers
- `producer/main.go`: enqueues immediate, delayed, and timeout-demo tasks
- `worker/main.go`: processes `critical` and `default` queues with weighted priority

## Run

Start worker in one terminal:

```bash
go run ./examples/basic/worker
```

To simulate abrupt worker crash behavior (without graceful shutdown on Ctrl+C):

```bash
go run ./examples/basic/worker -hard-exit-on-interrupt
```

Enqueue a task in another terminal:

```bash
go run ./examples/basic/producer
```

Expected behavior:

```text
Sending welcome email: user_id=42 email=user@example.com
Sending welcome email: user_id=7 email=delayed@example.com
Resizing image: source_url=https://example.com/assets/image.jpg
Resize canceled by context: source_url=https://example.com/assets/image.jpg err=context deadline exceeded
```

Notes:

- The delayed email task is enqueued with `ProcessIn(10 * time.Second)`.
- The image resize task is enqueued with `Timeout(5 * time.Second)`.
- The image resize handler simulates 8 seconds of work and returns `ctx.Err()` on cancellation, which triggers retries.
- `producer` supports selecting a single task type via `-task` flag: `all|email|delayed-email|image`.

## Crash failover demo

1. Start two workers in separate terminals:

```bash
go run ./examples/basic/worker -hard-exit-on-interrupt
go run ./examples/basic/worker -hard-exit-on-interrupt
```

2. Enqueue only image task:

```bash
go run ./examples/basic/producer -task image
```

3. Press Ctrl+C on the worker currently processing the image task.
4. Observe the other worker pick up the task after lease expiry and retry handoff.
