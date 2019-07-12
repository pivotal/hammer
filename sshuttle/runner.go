package sshuttle

import (
	"fmt"

	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (b Runner) Run(data environment.Config, dryRun bool, args ...string) error {
	networksPath := "/api/v0/staged/director/networks"
	cidrPath := ".networks[].subnets[].cidr"

	sshuttleCommandLines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`cidrs="$(om -t %s -k -u %s -p %s curl -s -p %s | jq -r %s | xargs echo)"`,
			data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, networksPath, cidrPath),

		fmt.Sprintf(`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@"%s" ${cidrs}`,
			data.OpsManager.IP.String()),
	}

	return b.ScriptRunner.RunScript(sshuttleCommandLines, []string{"jq", "om", "sshuttle"}, dryRun)
}
