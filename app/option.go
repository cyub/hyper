package app

// Options struct define app options
type Options struct {
	Name          string
	Addr          string
	RunMode       string
	CfgCenterAddr string
	CfgCenterPath string
	ShowBanner    bool
}

// Option use for inject option
type Option func(*Options)
