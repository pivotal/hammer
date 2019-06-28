package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/hammer/commands/commandsfakes"
	"github.com/pivotal/hammer/environment"

	. "github.com/pivotal/hammer/commands"
)

var _ = Describe("om command", func() {
	var (
		err     error
		command *OMCommand

		envReader *fakes.FakeEnvReader
		omRunner  *fakes.FakeToolRunner
		args      []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		omRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &OMCommand{
			Env:      envReader,
			OMRunner: omRunner,
			File:     true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the om tool", func() {
			Expect(omRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("runs the om tool using the retrieved environment config", func() {
			Expect(omRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := omRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(BeEquivalentTo([]string{"arg1", "arg2"}))
		})

		When("running the om tool is successful", func() {
			BeforeEach(func() {
				omRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the om tool errors", func() {
			BeforeEach(func() {
				omRunner.RunReturns(fmt.Errorf("om-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("om-runnner-error"))
			})
		})
	})
})
