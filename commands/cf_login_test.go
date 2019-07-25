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

var _ = Describe("cf login command", func() {
	var (
		err     error
		command *CFLoginCommand

		envReader     *fakes.FakeEnvReader
		ui            *fakes.FakeUI
		cfLoginRunner *fakes.FakeToolRunner
		commandArgs   []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		ui = new(fakes.FakeUI)
		cfLoginRunner = new(fakes.FakeToolRunner)
		commandArgs = []string{"arg1", "arg2"}

		command = &CFLoginCommand{
			Env:           envReader,
			UI:            ui,
			CFLoginRunner: cfLoginRunner,
			File:          true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(commandArgs)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the cf login tool", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			url, _ := url.Parse("www.test-cf.io")
			envReader.ReadReturns(environment.Config{OpsManager: environment.OpsManager{URL: *url}}, nil)
		})

		It("displays that the cf is being logged into", func() {
			Expect(ui.DisplayTextCallCount()).To(Equal(1))
			Expect(ui.DisplayTextArgsForCall(0)).To(Equal("Logging in to: www.test-cf.io\n"))
		})

		It("runs the cf login tool using the retrieved environment config", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

			environmentConfig, _, _ := cfLoginRunner.RunArgsForCall(0)

			expectedURL, _ := url.Parse("www.test-cf.io")
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{OpsManager: environment.OpsManager{URL: *expectedURL}}))
		})

		When("run with the file flag set", func() {
			BeforeEach(func() {
				command.File = true
			})

			It("runs the cf login tool in dry run mode", func() {
				Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := cfLoginRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeTrue())
			})
		})

		When("run with the file flag unset", func() {
			BeforeEach(func() {
				command.File = false
			})

			It("runs the cf login tool in non-dry run mode", func() {
				Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := cfLoginRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeFalse())
			})
		})

		It("runs the cf login tool with no additional args", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

			_, _, args := cfLoginRunner.RunArgsForCall(0)
			Expect(args).To(BeEmpty())
		})

		When("running the cf login tool is successful", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the cf login tool errors", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(fmt.Errorf("cf-login-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("cf-login-runnner-error"))
			})
		})
	})
})
