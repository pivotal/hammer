package commands_test

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/hammer/commands/commandsfakes"
	"github.com/pivotal/hammer/environment"

	. "github.com/pivotal/hammer/commands"
)

var _ = Describe("open command", func() {
	var (
		err     error
		command *OpenCommand

		envReader   *fakes.FakeEnvReader
		ui          *fakes.FakeUI
		openRunner  *fakes.FakeToolRunner
		commandArgs []string
	)

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		ui = new(fakes.FakeUI)
		openRunner = new(fakes.FakeToolRunner)
		commandArgs = []string{"arg1", "arg2"}

		command = &OpenCommand{
			Env:        envReader,
			UI:         ui,
			OpenRunner: openRunner,
			File:       true,
		}
	})

	JustBeforeEach(func() {
		err = command.Execute(commandArgs)
	})

	When("retrieving the environment config errors", func() {
		BeforeEach(func() {
			envReader.ReadReturns(environment.Config{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't attempt to run the open tool", func() {
			Expect(openRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("retrieving the environment config is successful", func() {
		BeforeEach(func() {
			url, _ := url.Parse("www.test-cf.io")
			envReader.ReadReturns(environment.Config{
				OpsManager: environment.OpsManager{
					URL:      *url,
					Username: "test-username",
					Password: "test-password"}},
				nil)
		})

		It("prints out that it is opening the ops manager and details on username and password", func() {
			Expect(ui.DisplayTextCallCount()).To(Equal(3))
			Expect(ui.DisplayTextArgsForCall(0)).To(Equal("Opening: www.test-cf.io\n"))
			Expect(ui.DisplayTextArgsForCall(1)).To(Equal("Username is: test-username\n"))
			Expect(ui.DisplayTextArgsForCall(2)).To(Equal("Password is in the clipboard\n"))
		})

		It("runs the open tool using the retrieved environment config", func() {
			Expect(openRunner.RunCallCount()).To(Equal(1))

			environmentConfig, _, _ := openRunner.RunArgsForCall(0)

			expectedUrl, _ := url.Parse("www.test-cf.io")
			Expect(environmentConfig).To(BeEquivalentTo(environment.Config{
				OpsManager: environment.OpsManager{
					URL:      *expectedUrl,
					Username: "test-username",
					Password: "test-password",
				},
			}))
		})

		When("run with the file flag set", func() {
			BeforeEach(func() {
				command.File = true
			})

			It("runs the open tool in dry run mode", func() {
				Expect(openRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := openRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeTrue())
			})
		})

		When("run with the file flag unset", func() {
			BeforeEach(func() {
				command.File = false
			})

			It("runs the open tool in non-dry run mode", func() {
				Expect(openRunner.RunCallCount()).To(Equal(1))

				_, dryRun, _ := openRunner.RunArgsForCall(0)
				Expect(dryRun).To(BeFalse())
			})
		})

		It("runs the open tool using the supplied command args", func() {
			Expect(openRunner.RunCallCount()).To(Equal(1))

			_, _, args := openRunner.RunArgsForCall(0)
			Expect(args).To(BeEmpty())
		})

		When("running the open tool is successful", func() {
			BeforeEach(func() {
				openRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("running the open tool errors", func() {
			BeforeEach(func() {
				openRunner.RunReturns(fmt.Errorf("open-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("open-runnner-error"))
			})
		})

		When("running in show mode", func() {
			BeforeEach(func() {
				command.Show = true
			})

			It("prints out the url, username and password for the ops manager", func() {
				Expect(ui.DisplayTextCallCount()).To(Equal(3))
				Expect(ui.DisplayTextArgsForCall(0)).To(Equal("www.test-cf.io\n"))
				Expect(ui.DisplayTextArgsForCall(1)).To(Equal("username: test-username\n"))
				Expect(ui.DisplayTextArgsForCall(2)).To(Equal("password: test-password\n"))
			})
		})
	})
})
