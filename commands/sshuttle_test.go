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
		args           []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		sshuttleRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &SshuttleCommand{
			Env:            envReader,
			SshuttleRunner: sshuttleRunner,
			File:           true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
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

			environmentConfig, dryRun, args := sshuttleRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
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
