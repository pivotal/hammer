/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package om

import (
	"fmt"
	"strings"

	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting"
)

type Runner struct {
	ScriptRunner scripting.ScriptRunner
}

func (r Runner) Run(data environment.Config, dryRun bool, omArgs ...string) error {
	var omCommandLines []string
	var omPrereqs []string

	if len(omArgs) > 0 {
		omCommandLines = []string{
			fmt.Sprintf(`om -t '%s' -k -u '%s' -p '%s' %s`, data.OpsManager.URL.String(), data.OpsManager.Username, data.OpsManager.Password, quoteArgs(omArgs)),
		}
		omPrereqs = []string{"om"}
	} else {
		omCommandLines = []string{
			fmt.Sprintf(`echo "export OM_TARGET=%s"`, data.OpsManager.URL.String()),
			fmt.Sprintf(`echo "export OM_USERNAME=%s"`, data.OpsManager.Username),
			fmt.Sprintf(`echo "export OM_PASSWORD=%s"`, data.OpsManager.Password),
		}
	}

	return r.ScriptRunner.RunScript(omCommandLines, omPrereqs, dryRun)
}

func quoteArgs(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, fmt.Sprintf(`'%s'`, arg))
	}
	return strings.Join(quoted, " ")
}
