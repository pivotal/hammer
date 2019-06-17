package open

import (
	"fmt"

	"github.com/pivotal/pcf-cli/scripting"

	"github.com/pivotal/pcf-cli/environment"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r Runner) Run(data environment.Config, dryRun bool, args ...string) error {
	openCommandLines := []string{
		fmt.Sprintf(`open "%s"`, data.OpsManager.URL.String()),
		fmt.Sprintf(`echo "%s" | pbcopy`, data.OpsManager.Password),
	}

	return r.ScriptRunner.RunScript(openCommandLines, []string{"open", "pbcopy"}, dryRun)
}
