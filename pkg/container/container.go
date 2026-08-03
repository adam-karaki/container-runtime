package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Config struct {
	Hostname string
	Command  []string
}

func getExecutablePath() (string, error) {
	if exe, err := os.Executable(); err == nil && exe != "" {
		return exe, nil
	}
	if _, err := os.Stat("/proc/self/exe"); err == nil {
		return "/proc/self/exe", nil
	}
	return "", fmt.Errorf("unable to resolve self executable path")
}

func RunParent(cfg Config) error {
	exePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	args := []string{"child"}
	if cfg.Hostname != "" {
		args = append(args, "-hostname", cfg.Hostname)
	}
	args = append(args, cfg.Command...)

	cmd := exec.Command(exePath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Delegates to OS-specific implementation
	cmd.SysProcAttr = getSysProcAttr()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start container process: %w", err)
	}

	return cmd.Wait()
}

func RunChild(hostname string, command []string) error {
	if err := setContainerHostname(hostname); err != nil {
		return err
	}

	binary, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("command not found %s: %w", command[0], err)
	}

	return syscall.Exec(binary, command, os.Environ())
}
