package common

import "time"

const (
	PersistenceFile = "./counter.json"
	ProbeInterval   = 1 * time.Second
	WindowSize      = 60 // 60 seconds window
)
