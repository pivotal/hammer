package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf-cli/commands/commandsfakes"
	"github.com/pivotal/pcf-cli/environment"

	. "github.com/pivotal/pcf-cli/commands"
)

var _ = Describe("Bosh command", func() {
	var (
		err     error
		command *BoshCommand

		envReader  *fakes.FakeEnvReader
		boshRunner *fakes.FakeToolRunner
		args       []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		boshRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &BoshCommand{
			Env:        envReader,
			BoshRunner: boshRunner,
			File:       true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the bosh tool", func() {
			Expect(boshRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("runs the bosh tool using the retrieved environment config", func() {
			Expect(boshRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := boshRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(BeEquivalentTo([]string{"arg1", "arg2"}))
		})

		When("running the bosh tool is successful", func() {
			BeforeEach(func() {
				boshRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the bosh tool errors", func() {
			BeforeEach(func() {
				boshRunner.RunReturns(fmt.Errorf("bosh-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("bosh-runnner-error"))
			})
		})
	})
})
