package om

import (
	"fmt"
	"strings"

	"github.com/pivotal/pcf/environment"
	"github.com/pivotal/pcf/scripting"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r Runner) Run(data environment.Config, dryRun bool, omArgs ...string) error {
	var omCommandLines []string

	if len(omArgs) > 0 {
		omCommandLines = []string{
			fmt.Sprintf(`om -t '%s' -k -u '%s' -p '%s' %s`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, quoteArgs(omArgs)),
		}
	} else {
		omCommandLines = []string{
			fmt.Sprintf(`echo "export OM_TARGET=%s"`, data.OpsManager.URL.String()),
			fmt.Sprintf(`echo "export OM_USERNAME=%s"`, data.OpsManager.Username),
			fmt.Sprintf(`echo "export OM_PASSWORD=%s"`, data.OpsManager.Password),
		}
	}

	return r.ScriptRunner.RunScript(omCommandLines, []string{"om"}, dryRun)
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, fmt.Sprintf(`'%s'`, arg))
	}
	return strings.Join(quoted, " ")
}
