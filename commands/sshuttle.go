/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands

type SshuttleCommand struct {
	TargetConfig    string `short:"t" long:"target" env:"HAMMER_TARGET_CONFIG" hidden:"true"`
	EnvironmentName string `short:"e" long:"environment-name" env:"HAMMER_ENVIRONMENT_NAME" hidden:"true"`
	File            bool   `short:"f" long:"file" description:"write a script file but do not run it"`

	Env            EnvReader
	UI             UI
	SshuttleRunner ToolRunner
}

func (c *SshuttleCommand) Execute(args []string) error {
	data, err := c.Env.Read(c.TargetConfig, c.EnvironmentName)
	if err != nil {
		return err
	}

	c.UI.DisplayText("# sshuttle\n")
	return c.SshuttleRunner.Run(data, c.File)
}
