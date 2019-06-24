package commands

import (
	"fmt"
)

type CFLoginCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env           EnvReader
	UI            UI
	CFLoginRunner ToolRunner
}

func (c *CFLoginCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	c.UI.DisplayText(fmt.Sprintf("Logging in to: %s\n", data.OpsManager.URL.String()))

	return c.CFLoginRunner.Run(data, c.File)
}
