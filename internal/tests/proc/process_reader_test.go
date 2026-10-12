package main

import (
	"pulsed/internal/proc"
	"testing"
)

func TestListPIDs(t *testing.T) {
	pids, err := proc.ListPIDs(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(pids) != 2 {
		t.Fatalf("expected 2PIDs, got %d", len(pids))
	}

	foundOne := false
	foundFortyTwo := false
	for _, pid := range pids {
		if pid == 1 {
			foundOne = true
		}
		if pid == 42 {
			foundFortyTwo = true
		}
	}
	if !foundOne {
		t.Fatalf("expected PID 1")
	}
	if !foundFortyTwo {
		t.Fatalf("expected PID 2")
	}

}
