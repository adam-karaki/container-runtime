package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// getExecutablePath resolves the binary path safely across Linux environments and test runners.
func getExecutablePath() (string, error) {
	if exe, err := os.Executable(); err == nil && exe != "" {
		return exe, nil
	}
	if _, err := os.Stat("/proc/self/exe"); err == nil {
		return "/proc/self/exe", nil
	}
	return "", fmt.Errorf("unable to resolve self executable path")
}

// RunParent spawns a child re-exec process of myrun that executes the target command.
func RunParent(command []string) error {
	exePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Re-exec current binary with the 'child' subcommand
	args := append([]string{"child"}, command...)

	cmd := exec.Command(exePath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Basic SysProcAttr configuration setup for execution
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start container process: %w", err)
	}

	return cmd.Wait()
}

// RunChild executes the actual target command inside the spawned process.
func RunChild(command []string) error {
	binary, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("command not found %s: %w", command[0], err)
	}

	// Replace the current child process image with the target command binary
	return syscall.Exec(binary, command, os.Environ())
}
