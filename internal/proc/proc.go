// SPDX-License-Identifier: MIT
package proc

import (
	"os"
	"strconv"
)

func ListPIDs(procDir string) ([]int, error) {
	entries, err := os.ReadDir(procDir)
	if err != nil {
		return nil, err
	}

	var pids []int

	for _, entry := range entries {
		if entry.IsDir() {
			if pid, err := strconv.Atoi(entry.Name()); err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}
