package cf

import (
	"fmt"

	"github.com/pivotal/pcf/scripting"

	"github.com/pivotal/pcf/environment"
)

type LoginRunner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r LoginRunner) Run(data environment.Config, dryRun bool, args ...string) error {
	lines := []string{
		fmt.Sprintf(`prods="$(om -t %s -k -u %s -p %s curl -s -p /api/v0/staged/products)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password),
		fmt.Sprintf(`guid="$(echo "$prods" | jq -r '.[] | select(.type == "cf") | .guid')"`),
		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password),
		fmt.Sprintf(`user="$(echo "$creds" | jq -r .credential.value.identity)"`),
		fmt.Sprintf(`pass="$(echo "$creds" | jq -r .credential.value.password)"`),
		fmt.Sprintf(`cf login -a "api.%s" -u "$user" -p "$pass" --skip-ssl-validation`, data.CFDomain),
	}

	return r.ScriptRunner.RunScript(lines, []string{"jq", "om", "cf"}, dryRun)
}
