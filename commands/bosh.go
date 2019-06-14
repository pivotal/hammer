package commands

type BoshCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env        EnvReader
	BoshRunner ToolRunner
}

func (c *BoshCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	return c.BoshRunner.Run(data, c.File, args...)
}
