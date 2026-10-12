package main

import (
	"os"
	"path/filepath"
	"pulsed/internal/proc"
	"testing"
)

func TestListPIDs(t *testing.T) {
	fixtureDir := t.TempDir()
	err := os.Mkdir(filepath.Join(fixtureDir, "1"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(fixtureDir, "42"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(fixtureDir, "notes"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	pids, err := proc.ListPIDs(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(pids) != 2 {
		t.Fatalf("expected 2 PIDs, got %d", len(pids))
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
		t.Fatalf("expected PID 42")
	}

}
