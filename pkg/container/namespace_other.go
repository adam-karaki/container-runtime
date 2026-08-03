//go:build !linux

package container

import (
	"fmt"
	"syscall"
)

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func setContainerHostname(hostname string) error {
	if hostname == "" {
		return nil
	}
	return fmt.Errorf("setting hostname requires Linux kernel primitives")
}
