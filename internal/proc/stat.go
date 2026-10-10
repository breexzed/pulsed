// SPDX-License-Identifier: MIT
package proc

import (
	"fmt"
	"pulsed/internal/model"
	"strings"
)

func readStat(pid int) (*model.ProcessInfo, error) {
	info := &model.ProcessInfo{PID: pid}

	return info, nil
}

func splitStat(line string) (comm string, fields []string, err error) {
	input := line
	open := strings.Index(input, "(")
	if open == -1 {
		return "", nil, fmt.Errorf("no open paren found in line")
	}
	close := strings.LastIndex(input, ")")
	if close == -1 {
		return "", nil, fmt.Errorf("no close paren found in line")
	}

	comm = input[open+1 : close]
	fields = strings.Fields(input[close+1:])

	return comm, fields, nil
}
