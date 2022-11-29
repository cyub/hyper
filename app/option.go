// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import "github.com/robfig/cron/v3"

// Options struct define app options
type Options struct {
	Name          string
	Addr          string
	RunMode       string
	CfgCenterAddr string
	CfgCenterPath string
	ShowBanner    bool
	CronEnable    bool
	Cron          *cron.Cron
}

// Option use for inject option
type Option func(*Options)
