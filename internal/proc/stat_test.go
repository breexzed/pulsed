package proc

import (
	"os"
	"testing"
)

func TestReadStatSelf(t *testing.T) {
	pid := os.Getpid()
	ppid := os.Getppid()

	procinfo, err := readStat(pid)
	if err != nil {
		t.Fatal(err)
	}
	if procinfo == nil {
		t.Fatal("expected process info")
	}
	if procinfo.PID != pid {
		t.Fatalf("expected PID %d, got %d", pid, procinfo.PID)
	}
	if procinfo.PPID != ppid {
		t.Fatalf("expected PPID %d, got %d", ppid, procinfo.PPID)
	}
	if procinfo.State == "" {
		t.Fatal("state not present")
	}
	if procinfo.Name == "" {
		t.Fatal("name not present")
	}
}

func TestSplitStat(t *testing.T) {
	line := "4242 (a b)) S 1 4242 4242 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"

	_, comm, fields, err := splitStat(line)
	if err != nil {
		t.Fatal(err)
	}
	if comm != "a b)" {
		t.Fatalf("expected comm %q, got %q", "a b)", comm)
	}
	if fields[0] != "S" {
		t.Fatalf("exected state %q, got %q", "1", fields[0])
	}
	if fields[1] != "1" {
		t.Fatalf("expected state %q, got %q", "1", fields[1])
	}
}

func TestParseStatFieldPositions(t *testing.T) {
	line := "1 (name) 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22"

	info, err := parseStat(line)
	if err != nil {
		t.Fatal(err)
	}
	if info.PID != 1 {
		t.Fatalf("expected PID %d, got %d", 1, info.PID)
	}
	if info.State != "3" {
		t.Fatalf("expected State 3, got %q", info.State)
	}
	if info.PPID != 4 {
		t.Fatalf("expected PPID %d, got %d", 4, info.PPID)
	}
	if info.UTime != 14 {
		t.Fatalf("expected UTime %d, got %d", 14, info.UTime)
	}
	if info.STime != 15 {
		t.Fatalf("expected STime %d, got %d", 15, info.STime)
	}
	if info.StartTime != 22 {
		t.Fatalf("expected StartTime %d, got %d", 22, info.StartTime)
	}
}
