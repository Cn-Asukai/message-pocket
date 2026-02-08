package cron

import (
	"context"
	"log/slog"

	"github.com/pocketbase/pocketbase/core"
	"github.com/samber/do/v2"
)

type Job struct {
	Name     string
	CronExpr string
	handle   func(ctx context.Context, i do.Injector) error
}

var jobs = make([]*Job, 0)

func Init(app core.App, i do.Injector) {
	for _, job := range jobs {
		err := app.Cron().Add(job.Name, job.CronExpr, func() {
			ctx := context.Background()
			err := job.handle(ctx, i)
			if err != nil {
				slog.ErrorContext(ctx, "job run fail", "name", job.Name, "err", err)
			}
		})
		if err != nil {
			panic(err)
		}
	}
}
