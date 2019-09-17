/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

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

		fmt.Sprintf(`ops_manager_ip="$(dig +short %s)"`, data.OpsManager.URL.Host),

		fmt.Sprintf(`cidrs="$(om -t %s -k -u %s -p %s curl -s -p %s | jq -r %s | xargs echo)"`,
			data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, networksPath, cidrPath),

		fmt.Sprintf(`sshuttle --ssh-cmd "ssh -o IdentitiesOnly=yes -i ${ssh_key_path}" -r ubuntu@${ops_manager_ip} ${cidrs}`),
	}

	return b.ScriptRunner.RunScript(sshuttleCommandLines, []string{"jq", "om", "sshuttle", "dig"}, dryRun)
}
