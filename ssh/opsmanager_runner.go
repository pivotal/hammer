/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package ssh

import (
	"fmt"

	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting"
)

type OpsManagerRunner struct {
	ScriptRunner scripting.ScriptRunner
}

func (b OpsManagerRunner) Run(data environment.Config, dryRun bool, args ...string) error {
	sshCommand := fmt.Sprintf(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" -t ubuntu@"%s"`, data.OpsManager.IP.String())

	sshCommandLines := []string{
		fmt.Sprintf(`ssh_key_path=$(mktemp)`),
		fmt.Sprintf(`echo "%s" >"$ssh_key_path"`, data.OpsManager.PrivateKey),
		fmt.Sprintf(`trap 'rm -f ${ssh_key_path}' EXIT`),
		fmt.Sprintf(`chmod 0600 "${ssh_key_path}"`),
		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p %s)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, "/api/v0/deployed/director/credentials/bosh_commandline_credentials"),
		fmt.Sprintf(`bosh="$(echo "$creds" | jq -r .credential | tr ' ' '\n' | grep '=')"`),
		fmt.Sprintf(`echo "$bosh"`),
		fmt.Sprintf(`shell="/usr/bin/env $(echo $bosh | tr '\n' ' ') bash -l"`),
		fmt.Sprintf(`%s "$shell"`, sshCommand),
	}

	return b.ScriptRunner.RunScript(sshCommandLines, []string{"ssh", "om"}, dryRun)
}
