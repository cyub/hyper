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
