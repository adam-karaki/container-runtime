package container

import (
	"fmt"
	"os"
	"os/exec"
)

type Config struct {
	Hostname string
	Rootfs   string
	Memory   string
	CPU      string
	Command  []string
}

func RunParent(cfg Config) error {
	args := []string{"child"}

	if cfg.Hostname != "" {
		args = append(args, "-hostname", cfg.Hostname)
	}
	if cfg.Rootfs != "" {
		args = append(args, "-rootfs", cfg.Rootfs)
	}

	args = append(args, cfg.Command...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = GetSysProcAttr()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start container process: %w", err)
	}

	// Apply cgroups to the child PID before waiting
	if err := ApplyCgroups(cfg, cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to apply cgroup limits: %w", err)
	}

	return cmd.Wait()
}

func RunChild(hostname string, rootfs string, args []string) error {
	if err := SetContainerHostname(hostname); err != nil {
		return err
	}

	if rootfs != "" {
		if err := PivotRoot(rootfs); err != nil {
			return err
		}
		if err := MountProc(); err != nil {
			return err
		}
	}

	if len(args) == 0 {
		return fmt.Errorf("no command provided to RunChild")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
