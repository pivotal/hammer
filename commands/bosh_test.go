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
	var err error
	var command *BoshCommand

	var envReader *fakes.FakeEnvReader
	var boshRunner *fakes.FakeToolRunner
	var args []string
	var dryRun bool

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		boshRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}
		dryRun = true
	})

	JustBeforeEach(func() {
		command = &BoshCommand{
			Env:        envReader,
			BoshRunner: boshRunner,
			File:       dryRun,
		}

		err = command.Execute(args)
	})

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call boshRunner", func() {
			Expect(boshRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{Name: "env-name"}, nil)
		})

		It("passes env data, args and dry run flag to boshRunner", func() {
			Expect(boshRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := boshRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(BeEquivalentTo([]string{"arg1", "arg2"}))
		})

		When("boshRunner succeeds", func() {
			BeforeEach(func() {
				boshRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("boshRunner returns an error", func() {
			BeforeEach(func() {
				boshRunner.RunReturns(fmt.Errorf("bosh-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("bosh-runnner-error"))
			})
		})
	})
})
