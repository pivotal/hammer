package commands

type OMCommand struct {
	Lockfile string `short:"l" long:"lockfile" env:"ENVIRONMENT_LOCK_METADATA" description:"path to a lockfile"`
	File     bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env      EnvReader
	OMRunner ToolRunner
}

func (c *OMCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.Lockfile)
	if err != nil {
		return err
	}

	return c.OMRunner.Run(data, c.File, args...)
}
