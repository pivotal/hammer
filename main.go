package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/pivotal/pcf/bosh"
	"github.com/pivotal/pcf/cf"
	"github.com/pivotal/pcf/commands"
	"github.com/pivotal/pcf/om"
	"github.com/pivotal/pcf/open"
	"github.com/pivotal/pcf/scripting"
	"github.com/pivotal/pcf/ssh"
	"github.com/pivotal/pcf/sshuttle"
)

var (
	version = "development build"
	date    = "unknown date"
)

type versionCommand struct{}

type lockfilePath struct{}

type options struct {
	Bosh       commands.BoshCommand       `command:"bosh" description:"display BOSH credentials, or run a BOSH command"`
	CFLogin    commands.CFLoginCommand    `command:"cf-login" description:"log in to cf on the environment"`
	Open       commands.OpenCommand       `command:"open" description:"open a browser to this environment"`
	OM         commands.OMCommand         `command:"om" description:"run the 'om' command with credentials for this environment"`
	SSH        commands.SSHCommand        `command:"ssh" description:"open an ssh connection to ops manager on this environment"`
	Sshuttle   commands.SshuttleCommand   `command:"sshuttle" description:"sshuttle to this environment"`
	Version    versionCommand             `command:"version" alias:"ver" description:"version of command"`
	Completion commands.CompletionCommand `command:"completion" description:"command completion script"`
	Lockfile   lockfilePath               `short:"l" long:"lockfile" hidden:"true"`
}

func main() {
	envReader := commands.NewEnvReader()
	scriptRunner := scripting.NewScriptRunner()

	opts := options{
		Bosh: commands.BoshCommand{
			Env: &envReader,
			BoshRunner: &bosh.Runner{
				ScriptRunner: scriptRunner,
			},
		},
		CFLogin: commands.CFLoginCommand{
			Env: &envReader,
			CFLoginRunner: &cf.LoginRunner{
				ScriptRunner: scriptRunner,
			},
		},
		OM: commands.OMCommand{
			Env: &envReader,
			OMRunner: &om.Runner{
				ScriptRunner: scriptRunner,
			},
		},
		Open: commands.OpenCommand{
			Env: &envReader,
			OpenRunner: &open.Runner{
				ScriptRunner: scriptRunner,
			},
		},
		SSH: commands.SSHCommand{
			Env: &envReader,
			SSHRunner: &ssh.Runner{
				ScriptRunner: scriptRunner,
			},
		},
		Sshuttle: commands.SshuttleCommand{
			Env: &envReader,
			SshuttleRunner: &sshuttle.Runner{
				ScriptRunner: scriptRunner,
			},
		},
	}

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
