package app

import (
	"github.com/robfig/cron/v3"
)

var (
	cronJobItems = make([]JobItem, 0)
)

// Job struct
type JobItem struct {
	Spec string
	cron.Job
}

// AddCronJob use for add job
func AddCronJob(spec string, job func()) {
	cronJobItems = append(cronJobItems, JobItem{spec, cron.FuncJob(job)})
}

// ApplyCron use for corn start
func applyCron(cron *cron.Cron) {
	if len(cronJobItems) == 0 {
		return
	}
	for _, item := range cronJobItems {
		cron.AddJob(item.Spec, item.Job)
	}
	cron.Start()
}
