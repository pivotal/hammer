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

			envReader *fakes.FakeEnvReader
			ui        *fakes.FakeUI
			sshRunner *fakes.FakeToolRunner
			args      []string
		)

		BeforeEach(func() {
			envReader = new(fakes.FakeEnvReader)
			ui = new(fakes.FakeUI)
			sshRunner = new(fakes.FakeToolRunner)
			args = []string{"arg1", "arg2"}

			sshOpsManagerCommand = &SSHOpsManagerCommand{
				Env:       envReader,
				UI:        ui,
				SSHRunner: sshRunner,
				File:      true,
			}
		})

		Context("no subcommand", func() {
			JustBeforeEach(func() {
				err = sshOpsManagerCommand.Execute(args)
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

					environmentConfig, dryRun, args := sshRunner.RunArgsForCall(0)
					Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
					Expect(dryRun).To(BeTrue())
					Expect(args).To(HaveLen(0))
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

	Context("ssh director subcommand", func() {
		var (
			err                error
			sshDirectorCommand *SSHDirectorCommand

			envReader *fakes.FakeEnvReader
			ui        *fakes.FakeUI
			sshRunner *fakes.FakeToolRunner
			args      []string
		)

		BeforeEach(func() {
			envReader = new(fakes.FakeEnvReader)
			ui = new(fakes.FakeUI)
			sshRunner = new(fakes.FakeToolRunner)
			args = []string{"arg1", "arg2"}

			sshDirectorCommand = &SSHDirectorCommand{
				Env:       envReader,
				UI:        ui,
				SSHRunner: sshRunner,
				File:      true,
			}
		})

		Context("no subcommand", func() {
			JustBeforeEach(func() {
				err = sshDirectorCommand.Execute(args)
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

					environmentConfig, dryRun, args := sshRunner.RunArgsForCall(0)
					Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
					Expect(dryRun).To(BeTrue())
					Expect(args).To(HaveLen(0))
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
})
