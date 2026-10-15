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
