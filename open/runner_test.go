package open_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/hammer/open"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/hammer/environment"
	"github.com/pivotal/hammer/scripting/scriptingfakes"
)

var _ = Describe("open runner", func() {
	var (
		err          error
		openRunner   open.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				URL:      *url,
				Password: "password",
			},
		}

		openRunner = open.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = openRunner.Run(data, dryRun)
	})

	It("runs the script with opsman url open and copying of the password into the clipboard", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, _, _ := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(Equal([]string{
			`open "www.test-url.io"`,
			`echo "password" | pbcopy`,
		}))
	})

	It("specifies the appropriate prerequisites when running the script", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		_, prereqs, _ := scriptRunner.RunScriptArgsForCall(0)

		Expect(prereqs).To(ConsistOf("open", "pbcopy"))
	})

	When("run with dry run set to false", func() {
		BeforeEach(func() {
			dryRun = false
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(false))
		})
	})

	When("run with dry run set to true", func() {
		BeforeEach(func() {
			dryRun = true
		})

		It("runs the script in dry run mode", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			_, _, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(dryRun).To(Equal(true))
		})
	})

	When("running the script succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
