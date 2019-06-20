package main

import (
	"fmt"
	"os"

	"github.com/pivotal/pcf-cli/environment"

	flags "github.com/jessevdk/go-flags"

	"github.com/pivotal/pcf-cli/bosh"
	"github.com/pivotal/pcf-cli/cf"
	"github.com/pivotal/pcf-cli/commands"
	"github.com/pivotal/pcf-cli/om"
	"github.com/pivotal/pcf-cli/open"
	"github.com/pivotal/pcf-cli/scripting"
	"github.com/pivotal/pcf-cli/ssh"
	"github.com/pivotal/pcf-cli/sshuttle"
)

var (
	version = "development build"
	date    = "unknown date"
)

type versionCommand struct{}

func (c *versionCommand) Execute(args []string) error {
	fmt.Printf("Version: %s (%s)\n", version, date)
	return nil
}

type targetConfigPath struct{}

// If `-t` is specified on the pcf command (rather than a subcommand)
// then set `TARGET_ENVIRONMENT_CONFIG` so the subcommand can read it
func (e *targetConfigPath) UnmarshalFlag(path string) error {
	return os.Setenv("TARGET_ENVIRONMENT_CONFIG", path)
}

type options struct {
	Bosh         commands.BoshCommand       `command:"bosh" description:"display BOSH credentials, or run a BOSH command"`
	CFLogin      commands.CFLoginCommand    `command:"cf-login" description:"log in to cf on the environment"`
	Open         commands.OpenCommand       `command:"open" description:"open a browser to this environment"`
	OM           commands.OMCommand         `command:"om" description:"run the 'om' command with credentials for this environment"`
	SSH          commands.SSHCommand        `command:"ssh" description:"open an ssh connection to ops manager on this environment"`
	Sshuttle     commands.SshuttleCommand   `command:"sshuttle" description:"sshuttle to this environment"`
	Version      versionCommand             `command:"version" alias:"ver" description:"version of command"`
	Completion   commands.CompletionCommand `command:"completion" description:"command completion script"`
	TargetConfig targetConfigPath           `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" description:"path to the target environment config"`
}

func main() {
	envReader := environment.Reader{}
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

func printDoubleDashMessage() {
	fmt.Fprintf(os.Stderr, "\nIf passing flags to 'bosh' or 'om', use a double dash '--', for example:\n")
	fmt.Fprintf(os.Stderr, "\n  pcf -t environment-path bosh -- -d deployment manifest\n")
}
