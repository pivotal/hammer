/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands_test

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/hammer/commands/commandsfakes"
	"github.com/pivotal/hammer/environment"

	. "github.com/pivotal/hammer/commands"
)

var _ = Describe("pks login command", func() {
	var (
		err     error
		command *PKSLoginCommand

		envReader      *fakes.FakeEnvReader
		ui             *fakes.FakeUI
		pksLoginRunner *fakes.FakeToolRunner
		commandArgs    []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		ui = new(fakes.FakeUI)
		pksLoginRunner = new(fakes.FakeToolRunner)
		commandArgs = []string{"arg1", "arg2"}

		command = &PKSLoginCommand{
			Env:            envReader,
			UI:             ui,
			PKSLoginRunner: pksLoginRunner,
			File:           true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(commandArgs)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the pks login tool", func() {
			Expect(pksLoginRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			url, _ := url.Parse("www.test-pks.io")
			envReader.ReadReturns(environment.Config{OpsManager: environment.OpsManager{URL: *url}}, nil)
		})

		It("displays that the pks is being logged into", func() {
			Expect(ui.DisplayTextCallCount()).To(Equal(2))
			Expect(ui.DisplayTextArgsForCall(0)).To(Equal("# pks-login\n"))
			Expect(ui.DisplayTextArgsForCall(1)).To(Equal("Logging in to PKS at: www.test-pks.io\n"))
		})

		It("runs the pks login tool using the retrieved environment config", func() {
			Expect(pksLoginRunner.RunCallCount()).To(Equal(1))

			environmentConfig, _, _ := pksLoginRunner.RunArgsForCall(0)

			expectedURL, _ := url.Parse("www.test-pks.io")
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{OpsManager: environment.OpsManager{URL: *expectedURL}}))
		})

		When("run with the file flag set", func() {
			BeforeEach(func() {
				command.File = true
			})

			It("runs the pks login tool in dry run mode", func() {
				Expect(pksLoginRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := pksLoginRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeTrue())
			})
		})

		When("run with the file flag unset", func() {
			BeforeEach(func() {
				command.File = false
			})

			It("runs the pks login tool in non-dry run mode", func() {
				Expect(pksLoginRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := pksLoginRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeFalse())
			})
		})

		It("runs the pks login tool with no additional args", func() {
			Expect(pksLoginRunner.RunCallCount()).To(Equal(1))

			_, _, args := pksLoginRunner.RunArgsForCall(0)
			Expect(args).To(BeEmpty())
		})

		When("running the pks login tool is successful", func() {
			BeforeEach(func() {
				pksLoginRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the pks login tool errors", func() {
			BeforeEach(func() {
				pksLoginRunner.RunReturns(fmt.Errorf("pks-login-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("pks-login-runnner-error"))
			})
		})
	})
})
