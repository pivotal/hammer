package actions

import (
	"fmt"
	"strings"

	"github.com/pivotal/pcf/lockfile"
)

type OMScripter struct{}

func NewOMScripter() OMScripter {
	return OMScripter{}
}

func (b OMScripter) Generate(data lockfile.Lockfile, omArgs []string) []string {
	omCommandLine := fmt.Sprintf(`om -t '%s' -k -u '%s' -p '%s' %s`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, quoteArgs(omArgs))

	return []string{omCommandLine}
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, fmt.Sprintf(`'%s'`, arg))
	}
	return strings.Join(quoted, " ")
}
