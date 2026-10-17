// SPDX-License-Identifier: MIT
package proc

import (
	"fmt"
	"os"
	"pulsed/internal/model"
	"strconv"
	"strings"
)

func parseStat(content string) (*model.ProcessInfo, error) {
	content = strings.TrimSpace(content)

	pid, comm, fields, err := splitStat(content)
	if err != nil {
		return nil, fmt.Errorf("splitting stat content: %w", err)
	}

	info := &model.ProcessInfo{PID: pid}
	info.Name = comm

	if len(fields) < 20 {
		return nil, fmt.Errorf("stat line has too few fields, got %d", len(fields))
	}

	state := fields[0]
	info.State = state

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, fmt.Errorf("parse ppid: %w", err)
	}
	info.PPID = ppid

	utime, err := strconv.ParseInt(fields[11], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse utime: %w", err)
	}
	info.UTime = utime

	stime, err := strconv.ParseInt(fields[12], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse stime: %w", err)
	}
	info.STime = stime

	startTime, err := strconv.ParseInt(fields[19], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse startTime: %w", err)
	}
	info.StartTime = startTime

	return info, nil
}

func readStat(pid int) (*model.ProcessInfo, error) {

	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, err
	}

	return parseStat(string(data))
}

func splitStat(line string) (pid int, comm string, fields []string, err error) {
	open := strings.Index(line, "(")
	if open == -1 {
		return 0, "", nil, fmt.Errorf("no open paren found in line")
	}
	closeParen := strings.LastIndex(line, ")")
	if closeParen == -1 {
		return 0, "", nil, fmt.Errorf("no close paren found in line")
	}
	pid, err = strconv.Atoi(strings.TrimSpace(line[:open]))
	if err != nil {
		return 0, "", nil, fmt.Errorf("split for pid: %w", err)
	}

	comm = line[open+1 : closeParen]
	fields = strings.Fields(line[closeParen+1:])

	return pid, comm, fields, nil
}
