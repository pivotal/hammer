/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package scripting

import (
	"fmt"
	"os"
	"os/exec"
)

//go:generate counterfeiter . ScriptRunner

type ScriptRunner interface {
	RunScript(lines []string, prereqs []string, onlyWriteFile bool) error
}

type scriptRunner struct{}

func NewScriptRunner() ScriptRunner {
	return scriptRunner{}
}

func (s scriptRunner) RunScript(lines []string, prereqs []string, onlyWriteFile bool) error {
	path, err := WriteTempFile(lines...)
	if err != nil {
		return err
	}

	if onlyWriteFile {
		fmt.Println(path)
		return nil
	}

	if err := CheckPrereqs(prereqs); err != nil {
		return err
	}

	defer os.Remove(path)

	debug := "+x"
	if len(os.Getenv("DEBUG")) > 0 {
		debug = "-x"
	}

	command := exec.Command("/bin/bash", debug, path)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}
