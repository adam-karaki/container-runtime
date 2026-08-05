package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/adam-karaki/container-runtime/pkg/container"
)

func Execute() error {
	if len(os.Args) < 2 {
		return errors.New("usage: myrun run [options] <command> [args...]")
	}

	subCmd := os.Args[1]
	switch subCmd {
	case "run":
		return parseAndRun(os.Args[2:])
	case "child":
		return parseAndRunChild(os.Args[2:])
	default:
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}
}

func parseAndRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	var cfg container.Config
	var memory, cpu, rootfs string

	fs.StringVar(&memory, "memory", "", "Limit memory usage")
	fs.StringVar(&cpu, "cpu", "", "Limit CPU usage")
	fs.StringVar(&cfg.Hostname, "hostname", "", "Container hostname")
	fs.StringVar(&rootfs, "rootfs", "", "Path to rootfs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return errors.New("must specify a command to run")
	}
	cfg.Command = cmdArgs

	return container.RunParent(cfg)
}

func parseAndRunChild(args []string) error {
	fs := flag.NewFlagSet("child", flag.ContinueOnError)

	var hostname, rootfs string
	fs.StringVar(&hostname, "hostname", "", "Container hostname")
	fs.StringVar(&rootfs, "rootfs", "", "Path to rootfs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return errors.New("child requires command execution arguments")
	}

	return container.RunChild(hostname, rootfs, cmdArgs)
}
