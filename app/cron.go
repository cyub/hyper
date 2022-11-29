// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
