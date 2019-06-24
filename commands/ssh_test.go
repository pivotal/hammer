package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf-cli/commands/commandsfakes"
	"github.com/pivotal/pcf-cli/environment"

	. "github.com/pivotal/pcf-cli/commands"
)

var _ = Describe("ssh command", func() {
	var (
		err     error
		command *SSHCommand

		envReader *fakes.FakeEnvReader
		sshRunner *fakes.FakeToolRunner
		args      []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		sshRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &SSHCommand{
			Env:       envReader,
			SSHRunner: sshRunner,
			File:      true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
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
