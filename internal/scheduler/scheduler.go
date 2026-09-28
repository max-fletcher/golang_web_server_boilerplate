package scheduler

import (
	"github.com/robfig/cron/v3"
)

// In order to add more CRON jobs, we need to create more methods that have a similar signature to "Add" method, the difference being
// the 2nd param. You may have a different function signature as 2nd param if the CRON job demands it. Then just register/start
// the method in server.go. That is it. Also, it is preferable to use the redis event queue in conjunction with these CRON jobs
// as it ensures that 2 don't run at the same time in case the job is not completed fast enough.

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (scheduler *Scheduler) Add(spec string, job func()) (cron.EntryID, error) {
	return scheduler.cron.AddFunc(spec, job)
}

func (scheduler *Scheduler) Start() {
	scheduler.cron.Start()
}

func (scheduler *Scheduler) Stop() {
	scheduler.cron.Stop()
}
