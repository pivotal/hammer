package open_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/pcf-cli/open"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf-cli/environment"
	"github.com/pivotal/pcf-cli/scripting/scriptingfakes"
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
		dryRun = true

		openRunner = open.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = openRunner.Run(data, dryRun)
	})

	Context("run", func() {
		It("invokes script runner with opsman url open and copying of the password into the clipboard", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ConsistOf(
				`open "www.test-url.io"`,
				`echo "password" | pbcopy`,
			))

			Expect(prereqs).To(ConsistOf("open", "pbcopy"))
			Expect(dryRun).To(Equal(true))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		When("script runner run script errors", func() {
			BeforeEach(func() {
				scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("run-script-error"))
			})
		})
	})
})
