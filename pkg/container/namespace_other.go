//go:build !linux

package container

import (
	"fmt"
	"syscall"
)

func GetSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func SetContainerHostname(hostname string) error {
	if hostname == "" {
		return nil
	}
	return fmt.Errorf("setting hostname is only supported on Linux")
}

func PivotRoot(rootfs string) error {
	if rootfs == "" {
		return nil
	}
	return fmt.Errorf("pivot_root is only supported on Linux")
}

func MountProc() error {
	return fmt.Errorf("mounting procfs is only supported on Linux")
}
