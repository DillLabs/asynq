# Asynq Scheduler Demo

This example demonstrates periodic task enqueueing via `asynq.NewScheduler` and task processing via a worker.

## Prerequisites

- Go installed
- Redis running on `127.0.0.1:6379`

## Files

- `tasks/tasks.go`: periodic task type, payload, and handler
- `worker/main.go`: worker process that consumes `report:generate` tasks
- `scheduler/main.go`: scheduler process that enqueues task every minute

## Run

Start worker in one terminal:

```bash
go run ./examples/scheduler/worker
```

Start scheduler in another terminal:

```bash
go run ./examples/scheduler/scheduler
```

## Expected behavior

- Scheduler logs registration info and enqueues one task every minute.
- Worker logs:

```text
Generating report: date=2026-04-27
```
