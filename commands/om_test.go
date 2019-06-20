package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf-cli/commands/commandsfakes"
	"github.com/pivotal/pcf-cli/environment"

	. "github.com/pivotal/pcf-cli/commands"
)

var _ = Describe("om command", func() {
	var err error
	var command *OMCommand

	var envReader *fakes.FakeEnvReader
	var omRunner *fakes.FakeToolRunner
	var args []string
	var dryRun bool

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

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call omRunner", func() {
			Expect(omRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("passes env data, args and dry run flag to omRunner", func() {
			Expect(omRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := omRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(BeEquivalentTo([]string{"arg1", "arg2"}))
		})

		When("omRunner succeeds", func() {
			BeforeEach(func() {
				omRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("omRunner returns an error", func() {
			BeforeEach(func() {
				omRunner.RunReturns(fmt.Errorf("om-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("om-runnner-error"))
			})
		})
	})
})
