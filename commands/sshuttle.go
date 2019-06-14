package commands

type SshuttleCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" description:"path to the target environment config"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env            EnvReader
	SshuttleRunner ToolRunner
}

func (c *SshuttleCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	return c.SshuttleRunner.Run(data, c.File)
}
