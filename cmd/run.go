package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/adam-karaki/container-runtime/pkg/container"
)

// Config holds CLI flags for myrun.
type Config struct {
	Memory   string
	CPU      string
	Hostname string
	RootFS   string
	Command  []string
}

// Execute parses command-line arguments and dispatches execution.
func Execute() error {
	if len(os.Args) < 2 {
		return errors.New("usage: myrun run [options] <command> [args...]")
	}

	subCmd := os.Args[1]
	switch subCmd {
	case "run":
		return parseAndRun(os.Args[2:])
	case "child":
		// Hidden re-exec entry point used for container child process execution
		return container.RunChild(os.Args[2:])
	default:
		return fmt.Errorf("unknown subcommand: %s", subCmd)
	}
}

func parseAndRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	var cfg Config
	fs.StringVar(&cfg.Memory, "memory", "", "Limit memory usage")
	fs.StringVar(&cfg.CPU, "cpu", "", "Limit CPU usage")
	fs.StringVar(&cfg.Hostname, "hostname", "", "Container hostname")
	fs.StringVar(&cfg.RootFS, "rootfs", "", "Path to rootfs")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return errors.New("must specify a command to run")
	}
	cfg.Command = cmdArgs

	return container.RunParent(cfg.Command)
}
