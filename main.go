package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/pivotal/hammer/bosh"
	"github.com/pivotal/hammer/cf"
	"github.com/pivotal/hammer/commands"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/om"
	"github.com/pivotal/hammer/open"
	"github.com/pivotal/hammer/scripting"
	"github.com/pivotal/hammer/ssh"
	"github.com/pivotal/hammer/sshuttle"
	"github.com/pivotal/hammer/ui"
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

// If `-t` is specified on the hammer command (rather than a subcommand)
// then set `TARGET_ENVIRONMENT_CONFIG` so the subcommand can read it
func (e *targetConfigPath) UnmarshalFlag(path string) error {
	return os.Setenv("TARGET_ENVIRONMENT_CONFIG", path)
}

type options struct {
	Bosh         commands.BoshCommand       `command:"bosh" description:"display BOSH credentials, or run a BOSH command"`
	CFLogin      commands.CFLoginCommand    `command:"cf-login" description:"log in to the cf for the environment"`
	Open         commands.OpenCommand       `command:"open" description:"open a browser to this environment"`
	OM           commands.OMCommand         `command:"om" description:"run the 'om' command with credentials for this environment"`
	SSH          commands.SSHCommand        `command:"ssh" choice:"opsman" choice:"director" description:"open an ssh connection to the ops manager or director of this environment"`
	Sshuttle     commands.SshuttleCommand   `command:"sshuttle" description:"sshuttle to this environment"`
	Version      versionCommand             `command:"version" alias:"ver" description:"version of command"`
	Completion   commands.CompletionCommand `command:"completion" description:"command completion script"`
	TargetConfig targetConfigPath           `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" description:"path to the target environment config"`
}

func main() {
	envReader := environment.Reader{}
	ui := ui.UI{
		Out: os.Stdout,
		Err: os.Stderr,
	}
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
			UI:  &ui,
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
			UI:  &ui,
			OpenRunner: &open.Runner{
				ScriptRunner: scriptRunner,
			},
		},
		SSH: commands.SSHCommand{
			OpsManager: commands.SSHOpsManagerCommand{
				Env: &envReader,
				UI:  &ui,
				SSHRunner: &ssh.OpsManagerRunner{
					ScriptRunner: scriptRunner,
				},
			},
			Director: commands.SSHDirectorCommand{
				Env: &envReader,
				UI:  &ui,
				SSHRunner: &ssh.DirectorRunner{
					ScriptRunner: scriptRunner,
				},
			},
		},
		Sshuttle: commands.SshuttleCommand{
			Env: &envReader,
			SshuttleRunner: &sshuttle.Runner{
				ScriptRunner: scriptRunner,
			},
		},
	}

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
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
	fmt.Fprintf(os.Stderr, "\n  hammer -t environment-path bosh -- -d deployment manifest\n")
}
