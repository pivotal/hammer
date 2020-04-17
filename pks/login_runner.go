/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package pks

import (
	"fmt"

	"github.com/pivotal/hammer/scripting"

	"github.com/pivotal/hammer/environment"
)

type LoginRunner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r LoginRunner) Run(data environment.Config, dryRun bool, args ...string) error {
	lines := []string{
		fmt.Sprintf(`prods="$(om -t %s -k -u %s -p %s curl -s -p /api/v0/staged/products)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password),
		`guid="$(echo "$prods" | jq -r '.[] | select(.type == "pivotal-container-service") | .guid')"`,
		fmt.Sprintf(`creds="$(om -t %s -k -u %s -p %s curl -s -p /api/v0/deployed/products/"$guid"/credentials/.properties.uaa_admin_password)"`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password),
		`pass="$(echo "$creds" | jq -r .credential.value.secret)"`,
		fmt.Sprintf(`pks login -a %s -u admin -p "$pass" --skip-ssl-validation`, data.PKSApi.URL.String()),
	}

	return r.ScriptRunner.RunScript(lines, []string{"jq", "om", "pks"}, dryRun)
}
