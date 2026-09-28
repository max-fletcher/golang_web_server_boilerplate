package jobs

import (
	"log/slog"
)

type TestJob struct {
	logger *slog.Logger
}

func NewTestJob(logger *slog.Logger) *TestJob {
	return &TestJob{
		logger: logger,
	}
}

func (job *TestJob) Run() {
	job.logger.Info("Test cron job executed")
}
