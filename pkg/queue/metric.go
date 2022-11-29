// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package queue

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	jobTotal       = "job_total"
	jobProcessTime = "job_process_time"
	successStatus  = "success"
	failStatus     = "fail"
)

var (
	metricOn      = false
	totalGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: jobTotal,
		Help: "Total number of jobs",
	}, []string{"name", "action", "status"})

	processingTimeVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    jobProcessTime,
		Help:    "The process time jobs in seconds",
		Buckets: nil,
	}, []string{"name"})
)

// MetricTurnon turn on metric exporter
func MetricTurnon() {
	metricOn = true
	prometheus.MustRegister(
		totalGaugeVec,
		processingTimeVec,
	)
}

// ConsumSuccessInc increment success consume job num
func ConsumSuccessInc(name string) {
	if !metricOn {
		return
	}
	totalGaugeVec.WithLabelValues(name, "consum", successStatus).Inc()
}

// ConsumFailInc increment failure consume job num
func ConsumFailInc(name string) {
	if !metricOn {
		return
	}
	totalGaugeVec.WithLabelValues(name, "consum", failStatus).Inc()
}

// InSuccessInc increment success enqueue job num
func InSuccessInc(name string) {
	if !metricOn {
		return
	}
	totalGaugeVec.WithLabelValues(name, "enqueue", successStatus).Inc()
}

// InFailInc incrment failure enqueue job num
func InFailInc(name string) {
	if !metricOn {
		return
	}
	totalGaugeVec.WithLabelValues(name, "enqueue", failStatus).Inc()
}

// ProcessTimeHist stats job process time
func ProcessTimeHist(name string, previous time.Time) {
	if !metricOn {
		return
	}
	processingTimeVec.WithLabelValues(name).Observe(time.Now().Sub(previous).Seconds())
}
