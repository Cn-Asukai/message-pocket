package cron

import (
	"context"
	"message-pocket/internal/services"

	"github.com/samber/do/v2"
)

func init() {
	jobs = append(jobs, &Job{
		Name:     "message_reprocessing",
		CronExpr: "",
		handle:   MessageReProcessing,
	})
}

func MessageReProcessing(ctx context.Context, i do.Injector) error {
	messageBoxService := do.MustInvoke[*services.MessageBoxService](i)
	err := messageBoxService.MessageRetry(ctx)
	if err != nil {
		return err
	}
	return nil
}
