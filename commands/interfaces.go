/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands

import "github.com/pivotal/hammer/environment"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . EnvReader

type EnvReader interface {
	Read(targetConfigPath string) (environment.Config, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ToolRunner

type ToolRunner interface {
	Run(data environment.Config, dryRun bool, args ...string) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . UI

type UI interface {
	DisplayText(text string)
	DisplayError(err error)
}
