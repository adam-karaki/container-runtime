//go:build linux

package container

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}
}

func setContainerHostname(hostname string) error {
	if hostname == "" {
		return nil
	}
	return syscall.Sethostname([]byte(hostname))
}
