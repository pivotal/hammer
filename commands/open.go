/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands

import (
	"fmt"
)

type OpenCommand struct {
	TargetConfig string `short:"t" long:"target" env:"TARGET_ENVIRONMENT_CONFIG" hidden:"true"`
	File         bool   `short:"f" long:"file" description:"write a script file but do not run it"`
	Show         bool   `short:"s" long:"show" description:"only show the credentials"`

	Env        EnvReader
	UI         UI
	OpenRunner ToolRunner
}

func (c *OpenCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig)
	if err != nil {
		return err
	}

	if c.Show {
		c.UI.DisplayText(fmt.Sprintf("%s\n", data.OpsManager.URL.String()))
		c.UI.DisplayText(fmt.Sprintf("username: %s\n", data.OpsManager.Username))
		c.UI.DisplayText(fmt.Sprintf("password: %s\n", data.OpsManager.Password))
		return nil
	}

	c.UI.DisplayText(fmt.Sprintf("Opening: %s\n", data.OpsManager.URL.String()))
	c.UI.DisplayText(fmt.Sprintf("Username is: %s\n", data.OpsManager.Username))
	c.UI.DisplayText(fmt.Sprint("Password is in the clipboard\n"))

	return c.OpenRunner.Run(data, c.File)
}
