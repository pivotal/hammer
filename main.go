package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal/pcf/internal/subcommands"
)

var (
	version = "development build"
	date    = "unknown date"
)

type versionCommand struct{}
type lockfilePath struct{}
type options struct {
	Bosh       subcommands.BoshCommand       `command:"bosh" description:"display BOSH credentials, or run a BOSH command"`
	CFLogin    subcommands.CFLoginCommand    `command:"cf-login" description:"log in to cf on the environment"`
	Open       subcommands.OpenCommand       `command:"open" description:"open a browser to this environment"`
	OM         subcommands.OMCommand         `command:"om" description:"run the 'om' command with credentials for this environment"`
	SSH        subcommands.SSHCommand        `command:"ssh" description:"open an ssh connection to ops manager on this environment"`
	Sshuttle   subcommands.SshuttleCommand   `command:"sshuttle" description:"sshutle to this environment"`
	Version    versionCommand                `command:"version" alias:"ver" description:"version of command"`
	Completion subcommands.CompletionCommand `command:"completion" description:"command completion script"`
	Lockfile   lockfilePath                  `short:"l" long:"lockfile" hidden:"true"`
}

func main() {
	var opts options

	if _, err := flags.Parse(&opts); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}

			if flagsErr.Type == flags.ErrUnknownFlag {
				printDoubleDashMessage()
			}
		}

		os.Exit(1)
	}
}

func (c *versionCommand) Execute(args []string) error {
	fmt.Printf("Version: %s (%s)\n", version, date)
	return nil
}

// If `-l` is specified on the pcf command (rather than a subcommand)
// then set `ENVIRONMENT_LOCK_METADATA` so the subcommand can read it
func (e *lockfilePath) UnmarshalFlag(path string) error {
	return os.Setenv("ENVIRONMENT_LOCK_METADATA", path)
}

func printDoubleDashMessage() {
	fmt.Fprintf(os.Stderr, "\nIf passing flags to 'bosh' or 'om', use a double dash '--', for example:\n")
	fmt.Fprintf(os.Stderr, "\n  pcf -l lockfile-path bosh -- -d deployment manifest\n")
}
