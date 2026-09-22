package smskit

import "time"

// Options contains shared client behavior that adapters may reuse.
type Options struct {
	UserAgent       string
	Timeout         time.Duration
	MaxResponseSize int64
}
