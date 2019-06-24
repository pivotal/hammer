package om_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/pcf-cli/om"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf-cli/environment"
	"github.com/pivotal/pcf-cli/scripting/scriptingfakes"
)

var _ = Describe("om runner", func() {
	var (
		err          error
		omRunner     om.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
		omArgs []string
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				URL:      *url,
				Username: "username",
				Password: "password",
			},
		}
		dryRun = true

		omRunner = om.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = omRunner.Run(data, dryRun, omArgs...)
	})

	When("no om args are passed to the runner", func() {
		BeforeEach(func() {
			omArgs = []string{}
		})

		It("runs the script with a series of om env var echos", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ConsistOf(
				`echo "export OM_TARGET=www.test-url.io"`,
				`echo "export OM_USERNAME=username"`,
				`echo "export OM_PASSWORD=password"`,
			))

			Expect(prereqs).To(HaveLen(0))
			Expect(dryRun).To(Equal(true))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("one or more om args are passed to the runner", func() {
		BeforeEach(func() {
			omArgs = []string{"arg1", "arg2", "arg3"}
			dryRun = false
		})

		It("runs the script with an om command", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ConsistOf(`om -t 'www.test-url.io' -k -u 'username' -p 'password' 'arg1' 'arg2' 'arg3'`))

			Expect(prereqs).To(ConsistOf("om"))
			Expect(dryRun).To(Equal(false))
		})

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
