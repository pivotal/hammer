package commands

type BoshCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" description:"path to the target environment config"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env        EnvReader
	BoshRunner ToolRunner
}

func (c *BoshCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	return c.BoshRunner.Run(data, c.File, args...)
}
