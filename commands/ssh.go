package commands

import (
	"fmt"
)

type SSHCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env       EnvReader
	UI        UI
	SSHRunner ToolRunner
}

func (c *SSHCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	c.UI.DisplayText(fmt.Sprintf("Connecting to: %s\n", data.Name))

	return c.SSHRunner.Run(data, c.File)
}
