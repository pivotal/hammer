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

var _ = Describe("ssh command", func() {
	Context("ssh opsman subcommand", func() {
		var (
			err                  error
			sshOpsManagerCommand *SSHOpsManagerCommand

			envReader   *fakes.FakeEnvReader
			ui          *fakes.FakeUI
			sshRunner   *fakes.FakeToolRunner
			commandArgs []string
		)

		BeforeEach(func() {
			envReader = new(fakes.FakeEnvReader)
			ui = new(fakes.FakeUI)
			sshRunner = new(fakes.FakeToolRunner)
			commandArgs = []string{"arg1", "arg2"}

			sshOpsManagerCommand = &SSHOpsManagerCommand{
				Env:       envReader,
				UI:        ui,
				SSHRunner: sshRunner,
				File:      true,
			}
		})

		JustBeforeEach(func() {
			err = sshOpsManagerCommand.Execute(commandArgs)
		})

		When("retrieving the environment config errors", func() {
			BeforeEach(func() {
				envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
			})

			It("doesn't attempt to run the ssh tool", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(0))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("env-reader-error"))
			})
		})

		When("retrieving the environment config is successful", func() {
			BeforeEach(func() {
				envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
			})

			It("displays that the connection is being started", func() {
				Expect(ui.DisplayTextCallCount()).To(Equal(1))
				Expect(ui.DisplayTextArgsForCall(0)).To(Equal("Connecting to: env-name\n"))
			})

			It("runs the ssh tool using the retrieved environment config", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(1))

				environmentConfig, _, _ := sshRunner.RunArgsForCall(0)
				Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			})

			When("run with the file flag set", func() {
				BeforeEach(func() {
					sshOpsManagerCommand.File = true
				})

				It("runs the ssh tool in dry run mode", func() {
					Expect(sshRunner.RunCallCount()).To(Equal(1))

					_, dryRun, _ := sshRunner.RunArgsForCall(0)
					Expect(dryRun).To(BeTrue())
				})
			})

			When("run with the file flag unset", func() {
				BeforeEach(func() {
					sshOpsManagerCommand.File = false
				})

				It("runs the ssh tool in non-dry run mode", func() {
					Expect(sshRunner.RunCallCount()).To(Equal(1))

					_, dryRun, _ := sshRunner.RunArgsForCall(0)
					Expect(dryRun).To(BeFalse())
				})
			})

			It("runs the ssh tool with no additional args", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(1))

				_, _, args := sshRunner.RunArgsForCall(0)
				Expect(args).To(BeEmpty())
			})

			When("running the ssh tool is successful", func() {
				BeforeEach(func() {
					sshRunner.RunReturns(nil)
				})

				It("doesn't error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("running the ssh tool errors", func() {
				BeforeEach(func() {
					sshRunner.RunReturns(fmt.Errorf("ssh-runnner-error"))
				})

				It("propagates the error", func() {
					Expect(err).To(MatchError("ssh-runnner-error"))
				})
			})
		})
	})

	Context("ssh director subcommand", func() {
		var (
			err                error
			sshDirectorCommand *SSHDirectorCommand

			envReader   *fakes.FakeEnvReader
			ui          *fakes.FakeUI
			sshRunner   *fakes.FakeToolRunner
			commandArgs []string
		)

		BeforeEach(func() {
			envReader = new(fakes.FakeEnvReader)
			ui = new(fakes.FakeUI)
			sshRunner = new(fakes.FakeToolRunner)
			commandArgs = []string{"arg1", "arg2"}

			sshDirectorCommand = &SSHDirectorCommand{
				Env:       envReader,
				UI:        ui,
				SSHRunner: sshRunner,
				File:      true,
			}
		})

		JustBeforeEach(func() {
			err = sshDirectorCommand.Execute(commandArgs)
		})

		When("retrieving the environment config errors", func() {
			BeforeEach(func() {
				envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
			})

			It("doesn't attempt to run the ssh tool", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(0))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("env-reader-error"))
			})
		})

		When("retrieving the environment config is successful", func() {
			BeforeEach(func() {
				envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
			})

			It("displays that the connection is being started", func() {
				Expect(ui.DisplayTextCallCount()).To(Equal(1))
				Expect(ui.DisplayTextArgsForCall(0)).To(Equal("Connecting to: env-name\n"))
			})

			It("runs the ssh tool using the retrieved environment config", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(1))

				environmentConfig, _, _ := sshRunner.RunArgsForCall(0)
				Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			})

			When("run with the file flag set", func() {
				BeforeEach(func() {
					sshDirectorCommand.File = true
				})

				It("runs the ssh tool in dry run mode", func() {
					Expect(sshRunner.RunCallCount()).To(Equal(1))

					_, dryRun, _ := sshRunner.RunArgsForCall(0)
					Expect(dryRun).To(BeTrue())
				})
			})

			When("run with the file flag unset", func() {
				BeforeEach(func() {
					sshDirectorCommand.File = false
				})

				It("runs the ssh tool in non-dry run mode", func() {
					Expect(sshRunner.RunCallCount()).To(Equal(1))

					_, dryRun, _ := sshRunner.RunArgsForCall(0)
					Expect(dryRun).To(BeFalse())
				})
			})

			It("runs the ssh tool with no additional args", func() {
				Expect(sshRunner.RunCallCount()).To(Equal(1))

				_, _, args := sshRunner.RunArgsForCall(0)
				Expect(args).To(BeEmpty())
			})

			When("running the ssh tool is successful", func() {
				BeforeEach(func() {
					sshRunner.RunReturns(nil)
				})

				It("doesn't error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("running the ssh tool errors", func() {
				BeforeEach(func() {
					sshRunner.RunReturns(fmt.Errorf("ssh-runnner-error"))
				})

				It("propagates the error", func() {
					Expect(err).To(MatchError("ssh-runnner-error"))
				})
			})
		})
	})
})
