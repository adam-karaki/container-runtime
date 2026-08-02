package container

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestBasicProcessExecution(t *testing.T) {
	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("failed to resolve test executable path: %v", err)
	}

	cmd := exec.Command(exePath, "child", "echo", "hello")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("expected clean exit, got: %v, stderr: %s", err, stderr.String())
	}

	got := strings.TrimSpace(stdout.String())
	want := "hello"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
