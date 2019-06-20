package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf-cli/commands/commandsfakes"
	"github.com/pivotal/pcf-cli/environment"

	. "github.com/pivotal/pcf-cli/commands"
)

var _ = Describe("cf login command", func() {
	var (
		err     error
		command *CFLoginCommand

		envReader     *fakes.FakeEnvReader
		cfLoginRunner *fakes.FakeToolRunner
		args          []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		cfLoginRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &CFLoginCommand{
			Env:           envReader,
			CFLoginRunner: cfLoginRunner,
			File:          true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
	})

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call cfLoginRunner", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("passes env data and dry run flag but not args to cfLoginRunner", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := cfLoginRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
		})

		When("cfLoginRunner succeeds", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("cfLoginRunner returns an error", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(fmt.Errorf("cf-login-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("cf-login-runnner-error"))
			})
		})
	})
})
