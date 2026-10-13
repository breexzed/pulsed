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
