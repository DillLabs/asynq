package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeWelcomeEmail = "email:welcome"
	TypeImageResize  = "image:resize"
)

type WelcomeEmailPayload struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type ImageResizePayload struct {
	SourceURL string `json:"source_url"`
}

func NewWelcomeEmailTask(userID int, email string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{
		UserID: userID,
		Email:  email,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageResize, payload), nil
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending welcome email: user_id=%d email=%s", p.UserID, p.Email)
	return nil
}

func HandleImageResizeTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: source_url=%s", p.SourceURL)

	// Simulate slow work and honor cancellation from task timeout.
	select {
	case <-time.After(8 * time.Second):
		log.Printf("Finished resize: source_url=%s", p.SourceURL)
		return nil
	case <-ctx.Done():
		log.Printf("Resize canceled by context: source_url=%s err=%v", p.SourceURL, ctx.Err())
		return ctx.Err()
	}
}
