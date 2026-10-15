package model

type ProcessInfo struct {
	PID       int
	Name      string
	State     string
	PPID      int
	Cmdline   []string
	UTime     int64 //cpu time spent in user mode - field[11]
	STime     int64 //cpu time spent in kernel mode - field[12]
	StartTime int64 //process start time since system booot - field[19]
}

// notes: User-mode time is spent executing application code,
//  while kernel-mode time is spent executing operating-system code on the process’s behalf.
