//go:build linux

package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func GetSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
	}
}

func SetContainerHostname(hostname string) error {
	if hostname == "" {
		return nil
	}
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return fmt.Errorf("failed to set hostname: %w", err)
	}
	return nil
}

func PivotRoot(rootfs string) error {
	if rootfs == "" {
		return nil
	}

	// Make mounts private so changes don't leak back to host
	if err := syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return fmt.Errorf("failed to make / private: %w", err)
	}

	// Bind mount rootfs onto itself (required by pivot_root)
	if err := syscall.Mount(rootfs, rootfs, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to bind-mount rootfs: %w", err)
	}

	pivotDir := filepath.Join(rootfs, ".pivot_old")
	if err := os.MkdirAll(pivotDir, 0700); err != nil {
		return fmt.Errorf("failed to create pivot_old: %w", err)
	}

	if err := syscall.PivotRoot(rootfs, pivotDir); err != nil {
		return fmt.Errorf("pivot_root failed: %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir to new root: %w", err)
	}

	oldPivotDir := filepath.Join("/", ".pivot_old")
	if err := syscall.Unmount(oldPivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("failed to unmount old root: %w", err)
	}

	return os.Remove(oldPivotDir)
}

func MountProc() error {
	_ = os.MkdirAll("/proc", 0755)
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("failed to mount /proc: %w", err)
	}
	return nil
}
func ApplyCgroups(cfg Config, pid int) error {
	if cfg.Memory == "" && cfg.CPU == "" {
		return nil
	}

	cgroupPath := "/sys/fs/cgroup/myrun"
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %w", err)
	}

	// Add process PID to cgroup
	procsPath := filepath.Join(cgroupPath, "cgroup.procs")
	if err := os.WriteFile(procsPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to attach pid to cgroup: %w", err)
	}

	// Memory Limit
	if cfg.Memory != "" {
		memPath := filepath.Join(cgroupPath, "memory.max")
		if err := os.WriteFile(memPath, []byte(cfg.Memory), 0644); err != nil {
			return fmt.Errorf("failed to write memory limit: %w", err)
		}
	}

	// CPU Limit (e.g. "0.5" -> "50000 100000")
	if cfg.CPU != "" {
		quota, err := parseCPUQuota(cfg.CPU)
		if err != nil {
			return err
		}
		cpuPath := filepath.Join(cgroupPath, "cpu.max")
		if err := os.WriteFile(cpuPath, []byte(quota), 0644); err != nil {
			return fmt.Errorf("failed to write cpu limit: %w", err)
		}
	}

	return nil
}

func parseCPUQuota(cpu string) (string, error) {
	val, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return "", fmt.Errorf("invalid cpu value %q: %w", cpu, err)
	}
	period := 100000
	quota := int(val * float64(period))
	return fmt.Sprintf("%d %d", quota, period), nil
}
