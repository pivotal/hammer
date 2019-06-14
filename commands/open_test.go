package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf/commands/commandsfakes"
	"github.com/pivotal/pcf/environment"

	. "github.com/pivotal/pcf/commands"
)

var _ = Describe("open command", func() {
	var err error
	var command *OpenCommand

	var envReader *fakes.FakeEnvReader
	var openRunner *fakes.FakeToolRunner
	var args []string
	var dryRun bool

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		openRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}
		dryRun = true
	})

	JustBeforeEach(func() {
		command = &OpenCommand{
			Env:        envReader,
			OpenRunner: openRunner,
			File:       dryRun,
		}

		err = command.Execute(args)
	})

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call openRunner", func() {
			Expect(openRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("passes env data and dry run flag but not args to openRunner", func() {
			Expect(openRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := openRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
		})

		When("openRunner succeeds", func() {
			BeforeEach(func() {
				openRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("openRunner returns an error", func() {
			BeforeEach(func() {
				openRunner.RunReturns(fmt.Errorf("open-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("open-runnner-error"))
			})
		})
	})
})
