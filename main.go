/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/pivotal/hammer/bosh"
	"github.com/pivotal/hammer/cf"
	"github.com/pivotal/hammer/commands"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/om"
	"github.com/pivotal/hammer/open"
	"github.com/pivotal/hammer/pks"
	"github.com/pivotal/hammer/scripting"
	"github.com/pivotal/hammer/ssh"
	"github.com/pivotal/hammer/sshuttle"
	"github.com/pivotal/hammer/ui"
)

var (
	version = "development build"
	date    = "unknown date"
)

type timeCommand struct{}

func (c *timeCommand) Execute([]string) error {
	fmt.Print("Can't touch this")
	return nil
}

type versionCommand struct{}

func (c *versionCommand) Execute([]string) error {
	fmt.Printf("Version: %s (%s)\n", version, date)
	return nil
}

type targetConfigPath struct{}

// If `-t` is specified on the hammer command (rather than a subcommand)
// then set `HAMMER_TARGET_CONFIG` so the subcommand can read it
func (e *targetConfigPath) UnmarshalFlag(path string) error {
	return os.Setenv("HAMMER_TARGET_CONFIG", path)
}

type options struct {
	Bosh         commands.BoshCommand       `command:"bosh" description:"display BOSH credentials, or run a BOSH command"`
	CFLogin      commands.CFLoginCommand    `command:"cf-login" description:"log in to the cf for the environment"`
	PKSLogin     commands.PKSLoginCommand   `command:"pks-login" description:"log in to pks for the environment"`
	Open         commands.OpenCommand       `command:"open" description:"open a browser to this environment"`
	OM           commands.OMCommand         `command:"om" description:"run the 'om' command with credentials for this environment"`
	SSH          commands.SSHCommand        `command:"ssh" choice:"opsman" choice:"director" description:"open an ssh connection to the ops manager or director of this environment"`
	Sshuttle     commands.SshuttleCommand   `command:"sshuttle" description:"sshuttle to this environment"`
	Time         timeCommand                `command:"time" description:"duuun dundundun" hidden:"true"`
	Version      versionCommand             `command:"version" alias:"ver" description:"version of command"`
	Completion   commands.CompletionCommand `command:"completion" description:"command completion script"`
	TargetConfig targetConfigPath           `short:"t" long:"target" env:"HAMMER_TARGET_CONFIG" description:"path to the target environment config"`
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
		PKSLogin: commands.PKSLoginCommand{
			Env: &envReader,
			UI:  &ui,
			PKSLoginRunner: &pks.LoginRunner{
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
		Time: timeCommand{},
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
