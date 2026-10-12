// SPDX-License-Identifier: MIT
package proc

import (
	"fmt"
	"os"
	"pulsed/internal/model"
	"strconv"
	"strings"
)

func readStat(pid int) (*model.ProcessInfo, error) {
	info := &model.ProcessInfo{PID: pid}

	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, err
	}

	content := strings.TrimSpace(string(data))

	comm, fields, err := splitStat(content)
	if err != nil {
		return nil, err
	}
	info.Name = comm

	if len(fields) < 2 {
		return nil, fmt.Errorf("stat line has too few fields")
	}

	state := fields[0]
	info.State = state

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}
	info.PPID = ppid

	return info, nil
}

func splitStat(line string) (comm string, fields []string, err error) {
	open := strings.Index(line, "(")
	if open == -1 {
		return "", nil, fmt.Errorf("no open paren found in line")
	}
	close := strings.LastIndex(line, ")")
	if close == -1 {
		return "", nil, fmt.Errorf("no close paren found in line")
	}

	comm = line[open+1 : close]
	fields = strings.Fields(line[close+1:])

	return comm, fields, nil
}
