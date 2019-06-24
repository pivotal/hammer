package commands_test

import (
	"fmt"
	"net/url"

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
		ui            *fakes.FakeUI
		cfLoginRunner *fakes.FakeToolRunner
		args          []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		ui = new(fakes.FakeUI)
		cfLoginRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}

		command = &CFLoginCommand{
			Env:           envReader,
			UI:            ui,
			CFLoginRunner: cfLoginRunner,
			File:          true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(args)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the cf login tool", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			url, _ := url.Parse("www.test-cf.io")
			envReader.ReadReturns(environment.Config{OpsManager: environment.OpsManager{URL: *url}}, nil)
		})

		It("displays that the cf is being logged into", func() {
			Expect(ui.DisplayTextCallCount()).To(Equal(1))
			Expect(ui.DisplayTextArgsForCall(0)).To(Equal("Logging in to: www.test-cf.io\n"))
		})

		It("runs the cf login tool using the retrieved environment config", func() {
			Expect(cfLoginRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := cfLoginRunner.RunArgsForCall(0)

			expectedUrl, _ := url.Parse("www.test-cf.io")
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{OpsManager: environment.OpsManager{URL: *expectedUrl}}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
		})

		When("running the cf login tool is successful", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the cf login tool errors", func() {
			BeforeEach(func() {
				cfLoginRunner.RunReturns(fmt.Errorf("cf-login-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("cf-login-runnner-error"))
			})
		})
	})
})
