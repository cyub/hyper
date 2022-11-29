// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// APIMetric is middleware use api metric
func APIMetric() gin.HandlerFunc {
	// the total http request
	HTTPReqTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	// the response time of http request
	HTTPReqDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP request latency in seconds",
		Buckets: nil,
	}, []string{"method", "path", "status"})

	// register metric
	prometheus.MustRegister(
		HTTPReqTotal,
		HTTPReqDuration,
	)

	// gin middleware
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		lvs := []string{
			c.Request.Method,
			c.Request.URL.Path,
			strconv.Itoa(c.Writer.Status()),
		}
		HTTPReqTotal.WithLabelValues(lvs...).Inc()
		HTTPReqDuration.WithLabelValues(lvs...).Observe(time.Now().Sub(startTime).Seconds())
	}
}
