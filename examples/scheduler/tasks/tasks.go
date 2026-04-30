package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/DillLabs/asynq"
)

const TypeReportGenerate = "report:generate"

type ReportPayload struct {
	Date string `json:"date"`
}

func NewReportTask(date string) (*asynq.Task, error) {
	payload, err := json.Marshal(ReportPayload{Date: date})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReportGenerate, payload), nil
}

func HandleReportTask(ctx context.Context, t *asynq.Task) error {
	var p ReportPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Generating report: date=%s", p.Date)
	return nil
}
