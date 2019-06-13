package commands

import (
	"fmt"
	"strings"

	"github.com/pivotal/pcf/internal/scripting"
)

type OMCommand struct {
	Env  EnvReader `group:"environment"`
	File bool      `short:"f" long:"file" description:"write a script file but do not run it"`
}

func (c *OMCommand) Execute(args []string) error {
	data, err := c.Env.Read()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		cmd := fmt.Sprintf(`om -t '%s' -k -u '%s' -p '%s' %s`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, quoteArgs(args))
		return scripting.RunScript([]string{cmd}, []string{"om"}, c.File)
	}

	fmt.Printf("export OM_TARGET=%s\n", data.OpsManager.URL.String())
	fmt.Printf("export OM_USERNAME=%s\n", data.OpsManager.Username)
	fmt.Printf("export OM_PASSWORD=%s\n", data.OpsManager.Password)
	return nil
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, fmt.Sprintf(`'%s'`, arg))
	}
	return strings.Join(quoted, " ")
}
