/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/hammer/commands/commandsfakes"
	"github.com/pivotal/hammer/environment"

	. "github.com/pivotal/hammer/commands"
)

var _ = Describe("sshuttle command", func() {
	var (
		err     error
		command *SshuttleCommand

		envReader      *fakes.FakeEnvReader
		sshuttleRunner *fakes.FakeToolRunner
		commandArgs    []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		sshuttleRunner = new(fakes.FakeToolRunner)
		commandArgs = []string{"arg1", "arg2"}

		command = &SshuttleCommand{
			Env:            envReader,
			SshuttleRunner: sshuttleRunner,
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

		It("doesn't attempt to run the sshuttle tool", func() {
			Expect(sshuttleRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("runs the sshuttle tool using the retrieved environment config", func() {
			Expect(sshuttleRunner.RunCallCount()).To(Equal(1))

			environmentConfig, _, _ := sshuttleRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
		})

		When("run with the file flag set", func() {
			BeforeEach(func() {
				command.File = true
			})

			It("runs the sshuttle tool in dry run mode", func() {
				Expect(sshuttleRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := sshuttleRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeTrue())
			})
		})

		When("run with the file flag unset", func() {
			BeforeEach(func() {
				command.File = false
			})

			It("runs the sshuttle tool in non-dry run mode", func() {
				Expect(sshuttleRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := sshuttleRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeFalse())
			})
		})

		It("runs the sshuttle tool using the supplied command args", func() {
			Expect(sshuttleRunner.RunCallCount()).To(Equal(1))

			_, _, args := sshuttleRunner.RunArgsForCall(0)
			Expect(args).To(BeEmpty())
		})

		When("running the sshuttle tool is successful", func() {
			BeforeEach(func() {
				sshuttleRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the sshuttle tool errors", func() {
			BeforeEach(func() {
				sshuttleRunner.RunReturns(fmt.Errorf("sshuttle-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("sshuttle-runnner-error"))
			})
		})
	})
})
