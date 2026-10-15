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
		return nil, err
	}

	info := &model.ProcessInfo{PID: pid}
	info.Name = comm

	if len(fields) < 20 {
		return nil, fmt.Errorf("stat line has too few fields")
	}

	state := fields[0]
	info.State = state

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}
	info.PPID = ppid

	Utime, err := strconv.Atoi(fields[11])
	if err != nil {
		return nil, err
	}
	info.UTime = int64(Utime)

	Stime, err := strconv.Atoi(fields[12])
	if err != nil {
		return nil, err
	}
	info.STime = int64(Stime)

	startTime, err := strconv.Atoi(fields[19])
	if err != nil {
		return nil, err
	}
	info.StartTime = int64(startTime)

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
	close := strings.LastIndex(line, ")")
	if close == -1 {
		return 0, "", nil, fmt.Errorf("no close paren found in line")
	}
	pid, err = strconv.Atoi(strings.TrimSpace(line[:open]))
	if err != nil {
		return 0, "", nil, err
	}

	comm = line[open+1 : close]
	fields = strings.Fields(line[close+1:])

	return pid, comm, fields, nil
}
