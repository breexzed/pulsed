package model

type ProcessInfo struct {
	PID     int
	Name    string
	State   string
	PPID    int
	Cmdline []string
}
