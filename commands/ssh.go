package commands

import (
	"fmt"
	"os"
)

type targetConfigPath struct{}

// If `-t` is specified on the ssh command (rather than a subcommand)
// then set `TARGET_ENVIRONMENT_CONFIG` so the subcommand can read it
func (e *targetConfigPath) UnmarshalFlag(path string) error {
	return os.Setenv("TARGET_ENVIRONMENT_CONFIG", path)
}

type SSHCommand struct {
	TargetConfig targetConfigPath     `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	Director     SSHDirectorCommand   `command:"director"`
	OpsManager   SSHOpsManagerCommand `command:"opsman"`
}

type SSHDirectorCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env       EnvReader
	UI        UI
	SSHRunner ToolRunner
}

func (c *SSHDirectorCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	c.UI.DisplayText(fmt.Sprintf("Connecting to: %s\n", data.Name))

	return c.SSHRunner.Run(data, c.File)
}

type SSHOpsManagerCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env       EnvReader
	UI        UI
	SSHRunner ToolRunner
}

func (c *SSHOpsManagerCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	c.UI.DisplayText(fmt.Sprintf("Connecting to: %s\n", data.Name))

	return c.SSHRunner.Run(data, c.File)
}
